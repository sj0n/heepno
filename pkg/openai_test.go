package pkg

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInit_AddsOpenAICmdToRoot(t *testing.T) {
	found := false
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "openai <file>" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("openaiCmd not registered with RootCmd")
	}
}

func TestOpenAICmd_FlagsDefaults(t *testing.T) {
	var oaiCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "openai <file>" {
			oaiCmd = cmd
			break
		}
	}
	if oaiCmd == nil {
		t.Fatalf("openaiCmd not found")
	}

	// Test translate flag default
	translate, err := oaiCmd.Flags().GetBool("translate")
	if err != nil {
		t.Fatalf("translate flag missing: %v", err)
	}
	if translate != false {
		t.Errorf("expected default translate to be false, got %v", translate)
	}

	// Test language flag default
	lang, err := oaiCmd.Flags().GetString("language")
	if err != nil {
		t.Fatalf("language flag missing: %v", err)
	}
	if lang != "" {
		t.Errorf("expected default language to be empty, got %q", lang)
	}

	// Test model flag default
	model, err := oaiCmd.Flags().GetString("model")
	if err != nil {
		t.Fatalf("model flag missing: %v", err)
	}
	if model != "whisper-1" {
		t.Errorf("expected default model to be 'whisper-1', got %q", model)
	}

	// Test format flag default
	format, err := oaiCmd.Flags().GetString("format")
	if err != nil {
		t.Fatalf("format flag missing: %v", err)
	}
	if format != "json" {
		t.Errorf("expected default format to be 'json', got %q", format)
	}

	// Test output flag default
	output, err := oaiCmd.Flags().GetString("output")
	if err != nil {
		t.Fatalf("output flag missing: %v", err)
	}
	if output != "" {
		t.Errorf("expected default output to be empty, got %q", output)
	}
}

func TestOpenAICmd_FlagsProperties(t *testing.T) {
	var oaiCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "openai <file>" {
			oaiCmd = cmd
			break
		}
	}
	if oaiCmd == nil {
		t.Fatalf("openaiCmd not found")
	}

	testCases := []struct {
		name      string
		shorthand string
		usage     string
	}{
		{"translate", "t", "Translate the audio file. Not setting this flag will transcribe the audio file."},
		{"language", "l", "Language of the source audio. Setting this helps in accuracy and velocity."},
		{"model", "m", "Model to use."},
		{"format", "f", "Format to use. json, text, srt, verbose_json, vtt"},
		{"output", "o", "The name of the output file. If not specified, the output will be printed to the console."},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag := oaiCmd.Flags().Lookup(tc.name)
			if flag == nil {
				t.Fatalf("flag %q not found", tc.name)
			}

			if flag.Shorthand != tc.shorthand {
				t.Errorf("expected shorthand for %q to be %q, got %q", tc.name, tc.shorthand, flag.Shorthand)
			}

			if flag.Usage != tc.usage {
				t.Errorf("expected usage for %q to be %q, got %q", tc.name, tc.usage, flag.Usage)
			}
		})
	}
}
