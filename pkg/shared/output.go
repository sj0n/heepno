package shared

import (
	"encoding/json"
	"fmt"
	"os"
)

// Print prints the data to the console in the specified format.
func Print(data any, text string, format string) error {
	switch format {
	case "json", "verbose_json":
		data, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("JSON Error: %v\n", err)
		}
		fmt.Println(string(data))
	default:
		fmt.Println(text)
	}

	return nil
}

// Save saves the data to the specified output file in the specified format.
func Save(data any, text string, format string, output string) error {
	cwd, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("File Error: %w", err)
	}

	switch format {
	case "json", "verbose_json":
		data, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("JSON Error: %w", err)
		}

		fileName, err := writeToFile(output, data, "json")
		if err != nil {
			return fmt.Errorf("File Error: %w", err)
		}

		fmt.Printf("Transcription saved to %s\\%s\n", cwd, fileName)
	default:
		fileName, err := writeToFile(output, text, "text")
		if err != nil {
			return fmt.Errorf("File Error: %w", err)
		}

		fmt.Printf("Transcription saved to %s\\%s\n", cwd, fileName)
	}

	return nil
}
