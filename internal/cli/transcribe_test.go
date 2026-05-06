package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sj0n/heepno/internal/provider"
)

func TestRequireAPIKey(t *testing.T) {
	testCases := []struct {
		name        string
		envKey      string
		envValue    string
		expectError bool
	}{
		{
			name:        "API key set",
			envKey:      "TEST_API_KEY",
			envValue:    "sk-test123",
			expectError: false,
		},
		{
			name:        "API key empty",
			envKey:      "TEST_API_KEY_EMPTY",
			envValue:    "",
			expectError: true,
		},
		{
			name:        "API key not set",
			envKey:      "TEST_API_KEY_NOT_SET",
			envValue:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				os.Setenv(tc.envKey, tc.envValue)
				defer os.Unsetenv(tc.envKey)
			} else {
				os.Unsetenv(tc.envKey)
			}

			err := RequireAPIKey(tc.envKey, "Test")
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "empty path",
			path:        "",
			expectError: true,
		},
		{
			name:        "file not found",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
		},
		{
			name:        "valid file",
			path:        createTestFile(t, "test"),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFile(tc.path)
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func createTestFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return path
}

func TestValidateFormat(t *testing.T) {
	allFormats := []string{"json", "text", "srt", "vtt", "verbose_json"}

	testCases := []struct {
		name        string
		format      string
		formats     []string
		expectError bool
	}{
		{
			name:        "JSON format",
			format:      "json",
			formats:     allFormats,
			expectError: false,
		},
		{
			name:        "TEXT format",
			format:      "text",
			formats:     allFormats,
			expectError: false,
		},
		{
			name:        "SRT format",
			format:      "srt",
			formats:     allFormats,
			expectError: false,
		},
		{
			name:        "VTT format",
			format:      "vtt",
			formats:     allFormats,
			expectError: false,
		},
		{
			name:        "VERBOSE_JSON format",
			format:      "verbose_json",
			formats:     allFormats,
			expectError: false,
		},
		{
			name:        "uppercase JSON",
			format:      "JSON",
			formats:     allFormats,
			expectError: false,
		},
		{
			name:        "invalid format",
			format:      "xml",
			formats:     allFormats,
			expectError: true,
		},
		{
			name:        "random format",
			format:      "csv",
			formats:     allFormats,
			expectError: true,
		},
		{
			name:        "srt not supported by provider",
			format:      "srt",
			formats:     []string{"json", "text"},
			expectError: true,
		},
		{
			name:        "vtt not supported by provider",
			format:      "vtt",
			formats:     []string{"json", "text"},
			expectError: true,
		},
		{
			name:        "json supported by provider",
			format:      "json",
			formats:     []string{"json", "text"},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFormat(tc.format, tc.formats)
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestValidateLanguage(t *testing.T) {
	testCases := []struct {
		name        string
		language    string
		expectError bool
	}{
		{
			name:        "empty language (allowed)",
			language:    "",
			expectError: false,
		},
		{
			name:        "two letter code",
			language:    "en",
			expectError: false,
		},
		{
			name:        "three letter code",
			language:    "eng",
			expectError: false,
		},
		{
			name:        "five letter code with hyphen",
			language:    "en-US",
			expectError: true,
		},
		{
			name:        "single letter",
			language:    "e",
			expectError: true,
		},
		{
			name:        "five letter code",
			language:    "english",
			expectError: true,
		},
		{
			name:        "code with numbers",
			language:    "en1",
			expectError: true,
		},
		{
			name:        "code with special chars",
			language:    "en!",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateLanguage(tc.language)
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestValidateModel(t *testing.T) {
	testCases := []struct {
		name        string
		provider    string
		model       string
		expectError bool
	}{
		{
			name:        "OpenAI whisper-1",
			provider:    "OpenAI",
			model:       "whisper-1",
			expectError: false,
		},
		{
			name:        "OpenAI uppercase",
			provider:    "OpenAI",
			model:       "WHISPER-1",
			expectError: false,
		},
		{
			name:        "OpenAI invalid model",
			provider:    "OpenAI",
			model:       "gpt-4",
			expectError: true,
		},
		{
			name:        "AssemblyAI universal",
			provider:    "AssemblyAI",
			model:       "universal",
			expectError: false,
		},
		{
			name:        "AssemblyAI slam-1",
			provider:    "AssemblyAI",
			model:       "slam-1",
			expectError: false,
		},
		{
			name:        "AssemblyAI invalid",
			provider:    "AssemblyAI",
			model:       "other",
			expectError: true,
		},
		{
			name:        "Deepgram any model",
			provider:    "Deepgram",
			model:       "nova-2",
			expectError: false,
		},
		{
			name:        "Deepgram custom model",
			provider:    "Deepgram",
			model:       "custom-model",
			expectError: false,
		},
		{
			name:        "Empty model",
			provider:    "OpenAI",
			model:       "",
			expectError: true,
		},
		{
			name:        "Unknown provider",
			provider:    "Other",
			model:       "any",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateModel(tc.provider, tc.model)
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestValidateOutputPath(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name        string
		outputPath  string
		expectError bool
	}{
		{
			name:        "empty output path (allowed)",
			outputPath:  "",
			expectError: false,
		},
		{
			name:        "valid output path",
			outputPath:  filepath.Join(tempDir, "output.txt"),
			expectError: false,
		},
		{
			name:        "non-existent directory",
			outputPath:  "/nonexistent/output.txt",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOutputPath(tc.outputPath)
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	openAISupported := []string{"json", "json_verbose", "text", "srt", "vtt"}
	testCases := []struct {
		name             string
		ctx              context.Context
		provider         string
		model            string
		lang             string
		format           string
		output           string
		supportedFormats []string
		transcribeFn     TranscribeFunc
		expectError      bool
	}{
		{
			name:     "success with console output",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en",
			format:   "text",
			output:   "",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{
					Text: "Hello world",
					Raw:  "transcript",
				}, nil
			},
			expectError: false,
		},
		{
			name:     "success with file output",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en",
			format:   "json",
			output:   filepath.Join(t.TempDir(), "output.json"),
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{
					Text: "Transcribed text",
					Raw:  map[string]string{"key": "value"},
				}, nil
			},
			expectError: false,
		},
		{
			name:     "empty text result",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en",
			format:   "text",
			output:   "",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{
					Text: "",
					Raw:  nil,
				}, nil
			},
			expectError: false,
		},
		{
			name:     "transcription error",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en",
			format:   "text",
			output:   "",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return nil, fmt.Errorf("transcription failed")
			},
			expectError: true,
		},
		{
			name:     "invalid format",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en",
			format:   "xml",
			output:   "",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{Text: "test", Raw: nil}, nil
			},
			expectError: true,
		},
		{
			name:     "invalid language",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en!",
			format:   "text",
			output:   "",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{Text: "test", Raw: nil}, nil
			},
			expectError: true,
		},
		{
			name:     "invalid model",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "invalid-model",
			lang:     "en",
			format:   "text",
			output:   "",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{Text: "test", Raw: nil}, nil
			},
			expectError: true,
		},
		{
			name:     "output path error",
			ctx:      context.Background(),
			provider: "OpenAI",
			model:    "whisper-1",
			lang:     "en",
			format:   "text",
			output:   "/nonexistent/output.txt",
			supportedFormats: openAISupported,
			transcribeFn: func(ctx context.Context) (*provider.Result, error) {
				return &provider.Result{Text: "test", Raw: nil}, nil
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Run(tc.ctx, tc.provider, tc.model, tc.lang, tc.format, tc.output, tc.supportedFormats, tc.transcribeFn)

			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestRun_NoSpeechDetected(t *testing.T) {
	transcribeFn := func(ctx context.Context) (*provider.Result, error) {
		return &provider.Result{
			Text: "",
			Raw:  nil,
		}, nil
	}

	err := Run(context.Background(), "OpenAI", "whisper-1", "en", "text", "", []string{"json", "text", "srt", "vtt", "verbose_json"}, transcribeFn)

	if err != nil {
		t.Errorf("expected no error for empty text, got: %v", err)
	}
}

func TestRun_TranscriptionTime(t *testing.T) {
	transcribeFn := func(ctx context.Context) (*provider.Result, error) {
		time.Sleep(50 * time.Millisecond)
		return &provider.Result{Text: "test", Raw: nil}, nil
	}

	start := time.Now()
	err := Run(context.Background(), "OpenAI", "whisper-1", "en", "text", "", []string{"json", "text", "srt", "vtt", "verbose_json"}, transcribeFn)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if duration > 150*time.Millisecond {
		t.Errorf("expected duration < 150ms, got %v", duration)
	}
}