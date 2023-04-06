package utils

import (
	"fmt"

    "github.com/fatih/color"
)

func Banner() {
    // Print a colorful welcome message
    fmt.Println()
    color.Blue("Welcome to Nucke Server!")
    fmt.Println()

    color.Yellow(`
        ,--,
  _ ___/ /\|
 ;( )__, )
; //   '--;
  \     |
   ^    ^`)

    fmt.Println()
}