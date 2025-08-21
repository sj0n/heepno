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

	testCases := []struct {
		name         string
		expected     any
		getFlagValue func() (any, error)
	}{
		{
			name:     "translate",
			expected: false,
			getFlagValue: func() (any, error) {
				return openaiCmd.Flags().GetBool("translate")
			},
		},
		{
			name:     "language",
			expected: "",
			getFlagValue: func() (any, error) {
				return openaiCmd.Flags().GetString("language")
			},
		},
		{
			name:     "format",
			expected: "json",
			getFlagValue: func() (any, error) {
				return oaiCmd.Flags().GetString("format")
			},
		},
		{
			name:     "output",
			expected: "",
			getFlagValue: func() (any, error) {
				return oaiCmd.Flags().GetString("output")
			},
		},
		{
			name:     "model",
			expected: "whisper-1",
			getFlagValue: func() (any, error) {
				return oaiCmd.Flags().GetString("model")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.getFlagValue()
			if err != nil {
				t.Errorf("error getting flag value: %v", err)
			}
			if value != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, value)
			}
		})
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
