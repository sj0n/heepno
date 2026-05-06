package cli

import (
	"testing"
)

func TestProviderConfig(t *testing.T) {
	testCases := []struct {
		name            string
		providerKey     string
		expectedName    string
		expectedAPIKey  string
		expectedModel   string
		expectTranslate bool
	}{
		{
			name:            "OpenAI provider",
			providerKey:     "OpenAI",
			expectedName:    "OpenAI",
			expectedAPIKey:  "OPENAI_API_KEY",
			expectedModel:   "whisper-1",
			expectTranslate: true,
		},
		{
			name:            "AssemblyAI provider",
			providerKey:     "AssemblyAI",
			expectedName:    "AssemblyAI",
			expectedAPIKey:  "ASSEMBLYAI_API_KEY",
			expectedModel:   "universal",
			expectTranslate: false,
		},
		{
			name:            "Deepgram provider",
			providerKey:     "Deepgram",
			expectedName:    "Deepgram",
			expectedAPIKey:  "DEEPGRAM_API_KEY",
			expectedModel:   "nova-3",
			expectTranslate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, ok := providers[tc.providerKey]
			if !ok {
				t.Fatalf("provider %q not found in providers map", tc.providerKey)
			}

			if cfg.Name != tc.expectedName {
				t.Errorf("expected Name %q, got %q", tc.expectedName, cfg.Name)
			}
			if cfg.APIKeyEnv != tc.expectedAPIKey {
				t.Errorf("expected APIKeyEnv %q, got %q", tc.expectedAPIKey, cfg.APIKeyEnv)
			}
			if cfg.ModelDefault != tc.expectedModel {
				t.Errorf("expected ModelDefault %q, got %q", tc.expectedModel, cfg.ModelDefault)
			}
			if cfg.SupportTranslate != tc.expectTranslate {
				t.Errorf("expected SupportTranslate %v, got %v", tc.expectTranslate, cfg.SupportTranslate)
			}
		})
	}
}

func TestProviderConfig_SupportedFormats(t *testing.T) {
	testCases := []struct {
		name            string
		providerKey     string
		expectedFormats []string
	}{
		{
			name:            "OpenAI formats",
			providerKey:     "OpenAI",
			expectedFormats: []string{"json", "json_verbose", "text", "srt", "vtt"},
		},
		{
			name:            "AssemblyAI formats",
			providerKey:     "AssemblyAI",
			expectedFormats: []string{"json", "text"},
		},
		{
			name:            "Deepgram formats",
			providerKey:     "Deepgram",
			expectedFormats: []string{"json", "text"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := providers[tc.providerKey]

			if len(cfg.SupportedFormats) != len(tc.expectedFormats) {
				t.Errorf("expected %d formats, got %d", len(tc.expectedFormats), len(cfg.SupportedFormats))
				return
			}

			for i, format := range tc.expectedFormats {
				if cfg.SupportedFormats[i] != format {
					t.Errorf("expected format[%d] %q, got %q", i, format, cfg.SupportedFormats[i])
				}
			}
		})
	}
}

func TestCreateCommand(t *testing.T) {
	testCases := []struct {
		name         string
		shortCommand string
		providerKey  string
	}{
		{
			name:         "OpenAI command",
			shortCommand: "oai",
			providerKey:  "OpenAI",
		},
		{
			name:         "Deepgram command",
			shortCommand: "dg",
			providerKey:  "Deepgram",
		},
		{
			name:         "AssemblyAI command",
			shortCommand: "aai",
			providerKey:  "AssemblyAI",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providerCfg := providers[tc.providerKey]
			cmd := createCommand(tc.shortCommand, providerCfg)

			if cmd == nil {
				t.Fatal("expected command, got nil")
			}

			expectedUse := tc.shortCommand + " <file>"
			if cmd.Use != expectedUse {
				t.Errorf("expected Use %q, got %q", expectedUse, cmd.Use)
			}

			expectedShort := "Transcribe audio using " + providerCfg.Name + "."
			if cmd.Short != expectedShort {
				t.Errorf("expected Short %q, got %q", expectedShort, cmd.Short)
			}

			if cmd.Args == nil {
				t.Error("expected Args validator, got nil")
			}
		})
	}
}

func TestCreateCommand_Flags(t *testing.T) {
	testCases := []struct {
		name                string
		providerKey         string
		expectTranslateFlag bool
		expectedFlags       []string
	}{
		{
			name:                "OpenAI flags",
			providerKey:         "OpenAI",
			expectTranslateFlag: true,
			expectedFlags:       []string{"language", "model", "format", "output", "translate"},
		},
		{
			name:                "Deepgram flags",
			providerKey:         "Deepgram",
			expectTranslateFlag: false,
			expectedFlags:       []string{"language", "model", "format", "output"},
		},
		{
			name:                "AssemblyAI flags",
			providerKey:         "AssemblyAI",
			expectTranslateFlag: false,
			expectedFlags:       []string{"language", "model", "format", "output"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providerCfg := providers[tc.providerKey]
			cmd := createCommand("test", providerCfg)

			flags := cmd.Flags()
			for _, flagName := range tc.expectedFlags {
				if flags.Lookup(flagName) == nil {
					t.Errorf("expected flag %q not found", flagName)
				}
			}
		})
	}
}

func TestCreateCommand_TranslateFlagOnlyForSupported(t *testing.T) {
	providersWithoutTranslate := []string{"Deepgram", "AssemblyAI"}

	for _, providerKey := range providersWithoutTranslate {
		t.Run(providerKey, func(t *testing.T) {
			providerCfg := providers[providerKey]
			cmd := createCommand("test", providerCfg)

			translateFlag := cmd.Flags().Lookup("translate")
			if providerCfg.SupportTranslate {
				if translateFlag == nil {
					t.Errorf("expected translate flag for %s", providerKey)
				}
			} else {
				if translateFlag != nil {
					t.Errorf("did not expect translate flag for %s", providerKey)
				}
			}
		})
	}
}
