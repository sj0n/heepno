package pkg

import (
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	Language        string
	OpenAIModel     string
	DeepgramModel   string
	AssemblyAIModel string
	Format          openai.AudioResponseFormat
	Translate       bool
	RootCmd         = &cobra.Command{
		Long: "A simple CLI tool to manage your tasks.",
		Version: "1.1.3",
	}
)