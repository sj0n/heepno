package main

import (
	"fmt"
	"os"

	"github.com/sj0n/heepno/pkg"
)

func main() {
	if err := pkg.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
