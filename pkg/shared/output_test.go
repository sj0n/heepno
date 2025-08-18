package shared

import (
	"bytes"
	"io"
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
			expectedOut: "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := Print(tc.data, tc.text, tc.format)

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = originalStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			gotOut := buf.String()

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
			}

			if gotOut != tc.expectedOut {
				t.Errorf("unexpected output.\ngot:\n%q\nwant:\n%q", gotOut, tc.expectedOut)
			}
		})
	}
}

func TestSave(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "heepno-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temp directory for the test
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
				content, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
					return
				}
				expected := "{\n  \"message\": \"test json\"\n}"
				if string(content) != expected {
					t.Errorf("unexpected content.\ngot:\n%s\nwant:\n%s", content, expected)
				}
			},
		},
		{
			name:        "Default text format",
			data:        nil,
			text:        "This is text content",
			format:      "text",
			output:      "test-text",
			expectError: false,
			extension:   ".txt",
			checkFunc: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
					return
				}
				expected := "This is text content"
				if string(content) != expected {
					t.Errorf("unexpected content.\ngot:\n%s\nwant:\n%s", content, expected)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout to ignore it in test output
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := Save(tc.data, tc.text, tc.format, tc.output)

			// Restore stdout
			w.Close()
			os.Stdout = originalStdout
			io.Copy(io.Discard, r)

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("did not expect an error but got: %v", err)
				return
			}

			// Check if file exists and has correct content
			filePath := tc.output + tc.extension
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("expected file %s was not created", filePath)
				return
			}

			tc.checkFunc(t, filePath)
		})
	}
}
