package initializers

import (
	"fmt"

    "github.com/fatih/color"
)

func Banner() {
    // Print a colorful welcome message
    white := color.New(color.FgWhite, color.Bold)

    white.Printf(`
           *                       *
                        .

                    _   __              __       
                   / | / /__  __ _____ / /__ ___ 
                  /  |/ // / / // ___// //_// _ \
                 / /|  // /_/ // /__ / ,<  /  __/
                /_/ |_/ \__,_/ \___//_/|_| \___/ 
                                                      
   *           
               
                                *
        |\___/|
        )     (             .              '
       =\     /=
         )===(       *
        /     \
        |     |
       /       \
       \       /
_/\_/\_/\__  _/_/\_/\_/\_/\_/\_/\_/\_/\_/\_
|  |  |  |( (  |  |  |  |  |  |  |  |  |  |
|  |  |  | ) ) |  |  |  |  |  |  |  |  |  |
|  |  |  |(_(  |  |  |  |  |  |  |  |  |  |
|  |  |  |  |  |  |  |  |  |  |  |  |  |  |
|  |  |  |  |  |  |  |  |  |  |  |  |  |  |
`)

    fmt.Println()
}