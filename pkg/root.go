package pkg

import (
	"github.com/spf13/cobra"
)

var (
	Language string
	Output   string
	Format   string
	RootCmd  = &cobra.Command{
		Long:    "Transcribe audio files using Deepgram, OpenAI and AssemblyAI models.",
		Version: "1.3.0",
	}
)
