package initializers

import (
	"fmt"

    "github.com/fatih/color"
)

func Banner() {
    // Print a colorful welcome message
    cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

    nucke := fmt.Sprint(`
                 _   __              __       
                / | / /__  __ _____ / /__ ___ 
               /  |/ // / / // ___// //_// _ \
              / /|  // /_/ // /__ / ,<  /  __/
             /_/ |_/ \__,_/ \___//_/|_| \___/ 
`)

    fmt.Printf(`
           *                       *
                        .

%s 
                                                      
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
`, cyan(nucke))

    fmt.Println()
}