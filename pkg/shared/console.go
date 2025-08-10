package shared

import (
	"fmt"
)

const (
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func PrintTranscriptionStatus(provider, model, language, status string) {
	fmt.Printf("%sProvider:%s  %s%s%s\n", ColorBold, ColorReset, ColorBlue, provider, ColorReset)
	fmt.Printf("%sModel:%s     %s%s%s\n", ColorBold, ColorReset, ColorBlue, model, ColorReset)
	fmt.Printf("%sLanguage:%s  %s%s%s\n", ColorBold, ColorReset, ColorBlue, language, ColorReset)
	fmt.Printf("%sStatus:%s    %s%s%s\n\n", ColorBold, ColorReset, ColorYellow, status, ColorReset)
}

// UpdateTranscriptionStatus replaces only the status line, keeping the rest unchanged.
// Call this after PrintTranscriptionStatus to update the status in place.
func UpdateTranscriptionStatus(status string, err error) {
	// Move cursor up 2 lines (from below the status line), return to start of line, and clear the line
	fmt.Print("\033[2A\r\033[K")

	if err != nil {
		fmt.Printf("%sStatus:%s    %s%s%s\n\n", ColorBold, ColorReset, ColorYellow, err, ColorReset)
		return
	}

	fmt.Printf("%sStatus:%s    %s%s%s\n\n", ColorBold, ColorReset, ColorYellow, status, ColorReset)
}
