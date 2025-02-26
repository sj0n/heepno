package pkg

import (
	"github.com/spf13/cobra"
)

var (
	Language string
	output   string
	format   string
	RootCmd  = &cobra.Command{
		Long:    "A simple CLI tool to manage your tasks.",
		Version: "1.2.0",
	}
)
