package console

import (
	"fmt"
	"testing"
)

func TestPrintTranscriptionStatus(t *testing.T) {
	testCases := []struct {
		name     string
		provider string
		model    string
		language string
		status   string
	}{
		{
			name:     "basic status",
			provider: "OpenAI",
			model:    "whisper-1",
			language: "en",
			status:   "Transcribing...",
		},
		{
			name:     "empty language",
			provider: "Deepgram",
			model:    "nova-2",
			language: "",
			status:   "Processing...",
		},
		{
			name:     "completed",
			provider: "AssemblyAI",
			model:    "universal",
			language: "es",
			status:   "Transcribed in 5s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PrintTranscriptionStatus panicked: %v", r)
				}
			}()
			PrintTranscriptionStatus(tc.provider, tc.model, tc.language, tc.status)
		})
	}
}

func TestUpdateTranscriptionStatus(t *testing.T) {
	testCases := []struct {
		name      string
		status    string
		err       error
		expectErr bool
	}{
		{
			name:      "success status",
			status:    "Transcribed in 3s",
			err:       nil,
			expectErr: false,
		},
		{
			name:      "error status",
			status:    "",
			err:       fmt.Errorf("transcription failed"),
			expectErr: true,
		},
		{
			name:      "empty status",
			status:    "no speech detected",
			err:       nil,
			expectErr: false,
		},
		{
			name:      "error message",
			status:    "",
			err:       fmt.Errorf("api error: rate limited"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("UpdateTranscriptionStatus panicked: %v", r)
				}
			}()
			UpdateTranscriptionStatus(tc.status, tc.err)
		})
	}
}
