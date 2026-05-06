package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Print prints the data to the console in the specified format.
func Print(data any, text string, format string) error {
	switch format {
	case "json", "verbose_json":
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("json error: %w", err)
		}
		fmt.Println(string(jsonBytes))
	default:
		fmt.Println(text)
	}

	return nil
}

// Save saves the data to the specified output file in the specified format.
func Save(data any, text string, format string, output string) error {
	cwd, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("file error: %w", err)
	}

	switch format {
	case "json", "verbose_json":
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("json error: %w", err)
		}

		fileName, err := writeToFile(output, jsonBytes, "json")
		if err != nil {
			return fmt.Errorf("file error: %w", err)
		}

		fmt.Printf("Transcription saved to %s\n", filepath.Join(cwd, fileName))
	case "text":
		fileName, err := writeToFile(output, text, "text")
		if err != nil {
			return fmt.Errorf("file error: %w", err)
		}
		fmt.Printf("Transcription saved to %s\n", filepath.Join(cwd, fileName))
	case "srt":
		fileName, err := writeToFile(output, text, "srt")
		if err != nil {
			return fmt.Errorf("file error: %w", err)
		}
		fmt.Printf("Transcription saved to %s\n", filepath.Join(cwd, fileName))
	case "vtt":
		fileName, err := writeToFile(output, text, "vtt")
		if err != nil {
			return fmt.Errorf("file error: %w", err)
		}
		fmt.Printf("Transcription saved to %s\n", filepath.Join(cwd, fileName))
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	return nil
}
