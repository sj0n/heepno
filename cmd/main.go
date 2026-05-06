package main

import (
	"fmt"
	"os"

	"github.com/sj0n/heepno/internal/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
