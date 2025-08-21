package pkg

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInit_AddsAAICmdToRoot(t *testing.T) {
	found := false
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "aai <file>" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("aaiCmd not registered with RootCmd")
	}
}

func TestAAICmd_FlagsDefaults(t *testing.T) {
	var aaiCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "aai <file>" {
			aaiCmd = cmd
			break
		}
	}
	if aaiCmd == nil {
		t.Fatalf("aaiCmd not found")
	}

	testCases := []struct {
		name         string
		expected     string
		getFlagValue func() (string, error)
	}{
		{
			name:     "language",
			expected: "",
			getFlagValue: func() (string, error) {
				return aaiCmd.Flags().GetString("language")
			},
		},
		{
			name:     "format",
			expected: "json",
			getFlagValue: func() (string, error) {
				return aaiCmd.Flags().GetString("format")
			},
		},
		{
			name:     "output",
			expected: "",
			getFlagValue: func() (string, error) {
				return aaiCmd.Flags().GetString("output")
			},
		},
		{
			name:     "model",
			expected: "universal",
			getFlagValue: func() (string, error) {
				return aaiCmd.Flags().GetString("model")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.getFlagValue()

			if err != nil {
				t.Errorf("error getting flag: %v", err)
			}

			if val != tc.expected {
				t.Errorf("Expecting %q got %q", tc.expected, val)
			}
		})
	}
}

func TestAAICmd_FlagsProperties(t *testing.T) {
	var aaiCmd *cobra.Command
	for _, cmd := range RootCmd.Commands() {
		if cmd.Use == "aai <file>" {
			aaiCmd = cmd
			break
		}
	}
	if aaiCmd == nil {
		t.Fatalf("aaiCmd not found")
	}

	testCases := []struct {
		name      string
		shorthand string
		usage     string
	}{
		{"language", "l", "Language to transcribe. See https://www.assemblyai.com/docs/getting-started/supported-languages for more details."},
		{"format", "f", "Transcribe format. <json|text>"},
		{"output", "o", "The name of the output file. If not specified, the output will be printed to the console."},
		{"model", "m", "Model to use. <universal|slam-1(only support English.)>"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag := aaiCmd.Flags().Lookup(tc.name)
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
