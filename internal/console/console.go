package console

import (
	"fmt"
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

func PrintTranscriptionStatus(provider, model, language, status string) {
	fmt.Printf("%sProvider:%s  %s%s%s\n", colorBold, colorReset, colorBlue, provider, colorReset)
	fmt.Printf("%sModel:%s     %s%s%s\n", colorBold, colorReset, colorBlue, model, colorReset)
	fmt.Printf("%sLanguage:%s  %s%s%s\n", colorBold, colorReset, colorBlue, language, colorReset)
	fmt.Printf("%sStatus:%s    %s%s%s\n\n", colorBold, colorReset, colorYellow, status, colorReset)
}

// UpdateTranscriptionStatus replaces only the status line, keeping the rest unchanged.
// Call this after PrintTranscriptionStatus to update the status in place.
func UpdateTranscriptionStatus(status string, err error) {
	// Move cursor up 2 lines (from below the status line), return to start of line, and clear the line
	fmt.Print("\033[2A\r\033[K")

	if err != nil {
		fmt.Printf("%sStatus:%s    %s%s%s\n\n", colorBold, colorReset, colorYellow, err, colorReset)
		return
	}

	fmt.Printf("%sStatus:%s    %s%s%s\n\n", colorBold, colorReset, colorYellow, status, colorReset)
}
