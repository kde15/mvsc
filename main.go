package main

import (
	"os"

	"github.com/kde15/mvsc/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
