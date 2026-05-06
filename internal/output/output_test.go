package output

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	type sampleData struct {
		Message string `json:"message"`
	}

	testCases := []struct {
		name        string
		data        any
		text        string
		format      string
		expectedOut string
		expectError bool
	}{
		{
			name:        "JSON format",
			data:        sampleData{Message: "hello"},
			text:        "",
			format:      "json",
			expectedOut: "{\n  \"message\": \"hello\"\n}\n",
			expectError: false,
		},
		{
			name:        "Verbose JSON format",
			data:        sampleData{Message: "hello verbose"},
			text:        "",
			format:      "verbose_json",
			expectedOut: "{\n  \"message\": \"hello verbose\"\n}\n",
			expectError: false,
		},
		{
			name:        "Default format (text)",
			data:        nil,
			text:        "simple text output",
			format:      "text",
			expectedOut: "simple text output\n",
			expectError: false,
		},
		{
			name:        "Unsupported format falls back to text",
			data:        nil,
			text:        "fallback text",
			format:      "srt",
			expectedOut: "fallback text\n",
			expectError: false,
		},
		{
			name:        "JSON marshal error",
			data:        make(chan int),
			text:        "some text",
			format:      "json",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Print(tc.data, tc.text, tc.format)

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
			}
		})
	}
}

func assertFileContent(t *testing.T, path, expected string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}
	if string(content) != expected {
		t.Errorf("unexpected content.\ngot:\n%s\nwant:\n%s", content, expected)
	}
}

func TestSave(t *testing.T) {
	tempDir := t.TempDir()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	type sampleData struct {
		Message string `json:"message"`
	}

	testCases := []struct {
		name        string
		data        any
		text        string
		format      string
		output      string
		expectError bool
		extension   string
		checkFunc   func(t *testing.T, path string)
	}{
		{
			name:        "JSON format",
			data:        sampleData{Message: "test json"},
			text:        "",
			format:      "json",
			output:      "test-json",
			expectError: false,
			extension:   ".json",
			checkFunc: func(t *testing.T, path string) {
				assertFileContent(t, path, "{\n  \"message\": \"test json\"\n}")
			},
		},
		{
			name:        "TEXT format",
			data:        nil,
			text:        "This is text content",
			format:      "text",
			output:      "test-text",
			expectError: false,
			extension:   ".txt",
			checkFunc: func(t *testing.T, path string) {
				assertFileContent(t, path, "This is text content")
			},
		},
		{
			name:        "SRT format",
			data:        nil,
			text:        "SRT content",
			format:      "srt",
			output:      "test-srt",
			expectError: false,
			extension:   ".srt",
			checkFunc: func(t *testing.T, path string) {
				assertFileContent(t, path, "SRT content")
			},
		},
		{
			name:        "VTT format",
			data:        nil,
			text:        "VTT content",
			format:      "vtt",
			output:      "test-vtt",
			expectError: false,
			extension:   ".vtt",
			checkFunc: func(t *testing.T, path string) {
				assertFileContent(t, path, "VTT content")
			},
		},
		{
			name:        "Unsupported format",
			data:        nil,
			text:        "Some text",
			format:      "xml",
			output:      "test-xml",
			expectError: true,
			extension:   "",
			checkFunc:   func(t *testing.T, path string) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Save(tc.data, tc.text, tc.format, tc.output)

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("did not expect an error but got: %v", err)
				return
			}

			filePath := tc.output + tc.extension
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("expected file %s was not created", filePath)
				return
			}

			tc.checkFunc(t, filePath)
		})
	}
}
