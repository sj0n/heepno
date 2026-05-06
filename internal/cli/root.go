package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Long:    "Transcribe audio files using Deepgram, OpenAI and AssemblyAI models.",
	Version: "1.9.1",
}

func init() {
	assemblyAICMD := createCommand("aai", providers["AssemblyAI"])
	deepgramCMD := createCommand("dg", providers["Deepgram"])
	openAICMD := createCommand("oai", providers["OpenAI"])

	RootCmd.AddCommand(assemblyAICMD)
	RootCmd.AddCommand(deepgramCMD)
	RootCmd.AddCommand(openAICMD)

	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true
}
