package cmd

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Commands that are used for authentication requests.",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
