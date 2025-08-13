package pkg

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Long:    "Transcribe audio files using Deepgram, OpenAI and AssemblyAI models.",
		Version: "1.4.0",
	}
)
