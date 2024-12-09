package pkg

import (
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	Language  string
	Model     string
	Format    openai.AudioResponseFormat
	Translate bool
	RootCmd   = &cobra.Command{
		Long: "A simple CLI tool to manage your tasks.",
	}
)
