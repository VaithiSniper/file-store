package util

import (
	"fmt"
	"strings"
)

const HyperstoreArt = `


 | |  | \ \   / /  __ \|  ____|  __ \ / ____|__   __/ __ \|  __ \|  ____|
 | |__| |\ \_/ /| |__) | |__  | |__) | (___    | | | |  | | |__) | |__   
 |  __  | \   / |  ___/|  __| |  _  / \___ \   | | | |  | |  _  /|  __|  
 | |  | |  | |  | |    | |____| | \ \ ____) |  | | | |__| | | \ \| |____ 
 |_|  |_|  |_|  |_|    |______|_|  \_\_____/   |_|  \____/|_|  \_\______|


`

type Color string

const (
	ColorGreen = "\u001b[32m"
	ColorBlue  = "\u001b[34m"
	ColorReset = "\u001b[0m"
)

func ColorPrint(color Color, message string) {
	fmt.Println(string(color), message, ColorReset)
}

const BannerWhiteSpaceChar = " "
const BannerBorderChar = "="
const BannerBorderCharVertical = "|"
const BannerBorderPadding = 4
const BannerBorderWidth = 60

func PrintInBanner(message string) {
	messageLength := len(message)
	computedBannerBorderWidth := BannerBorderWidth
	if messageLength+BannerBorderPadding > BannerBorderWidth {
		// Adjust border width if message is longer
		computedBannerBorderWidth = messageLength + BannerBorderPadding
	}

	border := strings.Repeat(BannerBorderChar, computedBannerBorderWidth)
	// Take away 3 for border chars and one whitespace at start
	whitespacePaddingLen := computedBannerBorderWidth - messageLength - 3
	whitespacePattern := strings.Repeat(
		BannerWhiteSpaceChar, whitespacePaddingLen,
	)
	fmt.Println(border)
	fmt.Printf(
		"%s%s%s%s%s\n", BannerBorderCharVertical, BannerWhiteSpaceChar, message,
		whitespacePattern,
		BannerBorderCharVertical,
	)
	fmt.Println(border)
}
