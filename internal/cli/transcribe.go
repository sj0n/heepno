package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/sj0n/heepno/internal/console"
	"github.com/sj0n/heepno/internal/output"
	"github.com/sj0n/heepno/internal/provider"
)

type TranscribeFunc func(context.Context) (*provider.Result, error)

func RequireAPIKey(key, name string) error {
	if os.Getenv(key) == "" {
		return fmt.Errorf("%s environment variable is not set", key)
	}
	return nil
}

// ValidateFile checks if the file exists and is readable
func ValidateFile(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	} else if err != nil {
		return fmt.Errorf("error accessing file: %w", err)
	}

	// Check if file is readable
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file is not readable: %w", err)
	}
	defer file.Close()

	return nil
}

// ValidateFormat checks if the format is supported
func ValidateFormat(format string) error {
	validFormats := []string{"json", "text", "srt", "vtt", "verbose_json"}
	for _, valid := range validFormats {
		if strings.ToLower(format) == valid {
			return nil
		}
	}
	return fmt.Errorf("unsupported format: %s (supported: json, text, srt, vtt, verbose_json)", format)
}

// ValidateLanguage checks if the language code is valid
func ValidateLanguage(lang string) error {
	if lang == "" {
		return nil // Empty language is allowed (use default)
	}

	// Basic validation: must be 2-4 characters
	if len(lang) < 2 || len(lang) > 4 {
		return fmt.Errorf("invalid language code length: %s (must be 2-4 characters)", lang)
	}

	// Check for valid characters (letters and optional hyphen)
	valid := true
	for _, c := range lang {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-') {
			valid = false
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid characters in language code: %s", lang)
	}

	return nil
}

// ValidateModel checks if the model is valid for the provider
func ValidateModel(providerName, model string) error {
	if model == "" {
		return fmt.Errorf("model cannot be empty")
	}

	// Provider-specific validations
	switch providerName {
	case "OpenAI":
		// OpenAI Whisper models
		validModels := []string{"whisper-1"}
		for _, valid := range validModels {
			if strings.ToLower(model) == valid {
				return nil
			}
		}
		return fmt.Errorf("unsupported OpenAI model: %s (supported: whisper-1)", model)

	case "AssemblyAI":
		// AssemblyAI models
		validModels := []string{"universal", "slam-1"}
		for _, valid := range validModels {
			if strings.ToLower(model) == valid {
				return nil
			}
		}
		return fmt.Errorf("unsupported AssemblyAI model: %s (supported: universal, slam-1)", model)

	case "Deepgram":
		// Deepgram models (allow any non-empty string as they add new models frequently)
		return nil

	default:
		// For unknown providers, just check non-empty
		return nil
	}
}

// ValidateOutputPath checks if the output path is valid
func ValidateOutputPath(path string) error {
	if path == "" {
		return nil // Empty output means print to console
	}

	// Check directory exists
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("output directory does not exist: %s", dir)
		} else if err != nil {
			return fmt.Errorf("error accessing output directory: %w", err)
		}
	}

	return nil
}

func Run(ctx context.Context, name, model, lang, format, out string, fn TranscribeFunc) error {
	// Validate all inputs before starting
	if err := ValidateFormat(format); err != nil {
		return err
	}

	if err := ValidateLanguage(lang); err != nil {
		return err
	}

	if err := ValidateModel(name, model); err != nil {
		return err
	}

	if err := ValidateOutputPath(out); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	console.PrintTranscriptionStatus(name, model, lang, "Transcribing...")

	start := time.Now()
	result, err := fn(ctx)
	if err != nil {
		return err
	}

	if result.Text == "" {
		console.UpdateTranscriptionStatus("no speech detected", nil)
		return nil
	}

	console.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", time.Since(start).Round(time.Second)), nil)

	if out != "" {
		return output.Save(result.Raw, result.Text, format, out)
	}
	return output.Print(result.Raw, result.Text, format)
}
