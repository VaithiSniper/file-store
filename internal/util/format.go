package util

import "fmt"

const HyperstoreArt = `


 | |  | \ \   / /  __ \|  ____|  __ \ / ____|__   __/ __ \|  __ \|  ____|
 | |__| |\ \_/ /| |__) | |__  | |__) | (___    | | | |  | | |__) | |__   
 |  __  | \   / |  ___/|  __| |  _  / \___ \   | | | |  | |  _  /|  __|  
 | |  | |  | |  | |    | |____| | \ \ ____) |  | | | |__| | | \ \| |____ 
 |_|  |_|  |_|  |_|    |______|_|  \_\_____/   |_|  \____/|_|  \_\______|


`

type Color string

const (
	ColorBlue  = "\u001b[34m"
	ColorReset = "\u001b[0m"
)

func ColorPrint(color Color, message string) {
	fmt.Println(string(color), message, ColorReset)
}
