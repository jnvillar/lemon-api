package cmd

import (
	"lemonapp/cmd/server"

	"github.com/spf13/cobra"
)

func Cmds() *cobra.Command {
	rootCmd := &cobra.Command{}

	runnable := server.NewRunnable()
	runnableCmd := runnable.Cmd()
	rootCmd.AddCommand(runnableCmd)

	return rootCmd
}
