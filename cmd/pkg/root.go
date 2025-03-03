package pkg

import (
	"github.com/spf13/cobra"
)

var (
	Language string
	output   string
	format   string
	RootCmd  = &cobra.Command{
		Long:    "Transcribe audio files using Deepgram, OpenAI and SeemblyAI models.",
		Version: "1.2.1",
	}
)
