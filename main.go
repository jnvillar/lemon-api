package main

import (
	"os"

	"lemonapp/cmd"
)

func main() {
	if err := cmd.Cmds().Execute(); err != nil {
		os.Exit(1)
	}
}
