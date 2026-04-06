package console

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintTranscriptionStatus(provider, model, language, status string) {
	bold := color.New(color.Bold)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	bold.Println("Provider: ", blue.Sprint(provider))
	bold.Println("Model:     ", blue.Sprint(model))
	bold.Println("Language:  ", blue.Sprint(language))
	bold.Println("Status:    ", yellow.Sprint(status))
	fmt.Println()
}

func UpdateTranscriptionStatus(status string, err error) {
	fmt.Print("\033[2A\r\033[K")

	bold := color.New(color.Bold)
	yellow := color.New(color.FgYellow)

	if err != nil {
		bold.Println("Status:    ", yellow.Sprint(err.Error()))
		fmt.Println()
		return
	}

	bold.Println("Status:    ", yellow.Sprint(status))
	fmt.Println()
}
