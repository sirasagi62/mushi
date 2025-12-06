package main

import (
	"os"

	"github.com/sirasagi62/mushi/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
