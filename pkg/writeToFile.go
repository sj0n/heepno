package pkg

import (
	"fmt"
	"os"
)

func writeToFile(filename string, data any, format string) (string, error) {
	var ext string

	switch format {
	case "json", "verbose_json":
		ext = ".json"
	case "text":
		ext = ".txt"
	case "srt":
		ext = ".srt"
	case "vtt":
		ext = ".vtt"
	}

	fmt.Println("Saving to file...")

	file, err := os.Create(filename + ext)
	if err != nil {
		return "", fmt.Errorf("file error %w", err)
	}
	defer file.Close()

	switch content := data.(type) {
	case string:
		_, err = file.WriteString(content)

		if err != nil {
			return "", fmt.Errorf("write error %w", err)
		}
	case []byte:
		_, err = file.Write(content)

		if err != nil {
			return "", fmt.Errorf("write error %w", err)
		}
	}

	return file.Name(), nil
}
