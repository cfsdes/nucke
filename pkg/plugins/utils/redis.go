package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cfsdes/nucke/pkg/report"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
)

func RunRedis() {
	Red := color.New(color.FgRed, color.Bold).SprintFunc()
	Cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("[%s] Starting Redis...\n", Cyan("INF"))

	// Inicia o Redis
	cmd := exec.Command("redis-server")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("[%s] Error while initializing redis: %v\n", Red("ERR"), err)
		return
	}

	time.Sleep(2 * time.Second)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Checar por interações passadas do interactsh
	checkOobInteraction(client)

	// Limpar registros antigos
	ctx := context.Background()
	client.FlushAll(ctx)

	// Monitora alterações no arquivo /tmp/nucke-interact
	go watchFile("/tmp/nucke-interact", func() {
		checkOobInteraction(client)
	})
}

// Armazena as infos do scan OOB no Redis
func StoreDetection(pluginDir, oob_id, url, payload, param, rawReq string) {
	scanName := filepath.Base(pluginDir)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	session := map[string]string{
		"scanName":  scanName,
		"url":       url,
		"payload":   payload,
		"param":     param,
		"rawReq":    rawReq,
		"pluginDir": pluginDir,
	}

	for k, v := range session {
		err := client.HSet(ctx, oob_id, k, v).Err()
		if err != nil {
			panic(err)
		}
	}

	// Set expiration time - 6 hours
	expTime := time.Duration(6) * time.Hour
	client.Expire(ctx, oob_id, expTime)
}

func watchFile(caminho string, callback func()) {
	Red := color.New(color.FgRed, color.Bold).SprintFunc()

	// Criação do watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("[%s] Watcher error: %v\n", Red("ERR"), err)
		return
	}
	defer watcher.Close()

	// Adicionar arquivo ao watcher
	err = watcher.Add(caminho)
	if err != nil {
		fmt.Printf("[%s] Watcher error: %v\n", Red("ERR"), err)
		return
	}

	// Goroutine para processar eventos
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					// O arquivo foi alterado, chama o callback
					callback()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("[%s] Watcher error: %v\n", Red("ERR"), err)
			}
		}
	}()

	// Aguardar indefinidamente para manter o programa em execução
	select {}
}

func get8digitsKey(client *redis.Client) ([]string, error) {
	// Obter todas as chaves do Redis
	ctx := context.Background()
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	// Filtrar chaves com exatamente 8 dígitos
	var chavesCom8Digitos []string
	for _, key := range keys {
		if has8digits(key) {
			chavesCom8Digitos = append(chavesCom8Digitos, key)
		}
	}

	return chavesCom8Digitos, nil
}

func has8digits(s string) bool {
	// Utiliza uma expressão regular para verificar se a string possui exatamente 8 dígitos
	match, _ := regexp.MatchString("^[0-9]{8}$", s)
	return match
}

func verifyKeys(chaves []string) (bool, string) {
	Red := color.New(color.FgRed, color.Bold).SprintFunc()

	// Abre o arquivo para leitura
	arquivo, err := os.Open("/tmp/nucke-interact")
	if err != nil {
		fmt.Printf("[%s] Watcher error: %v\n", Red("ERR"), err)
		return false, ""
	}
	defer arquivo.Close()

	// Cria um scanner para ler as linhas do arquivo
	scanner := bufio.NewScanner(arquivo)

	// Percorre cada linha do arquivo
	for scanner.Scan() {
		linha := scanner.Text()

		// Verifica se cada chave do Redis está presente na linha do arquivo
		for _, chave := range chaves {
			if strings.Contains(linha, chave) {
				reportFinding(chave)
			}
		}
	}

	// Verifica se ocorreu algum erro durante a leitura do arquivo
	if err := scanner.Err(); err != nil {
		fmt.Printf("[%s] Watcher error: %v\n", Red("ERR"), err)
	}

	return false, ""
}

func reportFinding(key string) {
	// Redis Client
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Report finding
	vuln := client.HGetAll(ctx, key).Val()
	webhook := report.GetWebhook(vuln["pluginDir"] + "/" + vuln["scanName"] + ".so")
	report.Output(vuln["scanName"], webhook, "OOB", vuln["url"], vuln["payload"], vuln["param"], vuln["rawReq"], "", vuln["pluginDir"])

	// Clear key in redis
	client.Del(ctx, key)
}

// Checar se o interactsh teve interação
func checkOobInteraction(client *redis.Client) {
	Red := color.New(color.FgRed, color.Bold).SprintFunc()

	// Obter as chaves do Redis com 8 dígitos
	chaves, err := get8digitsKey(client)
	if err != nil {
		fmt.Printf("[%s] Error getting redis key: %v\n", Red("ERR"), err)
		return
	}

	// Verificar se as chaves do Redis existem na sessao do interactsh
	verifyKeys(chaves)
}
