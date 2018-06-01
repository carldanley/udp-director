package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run:   runVersionCommand,
	}

	RootCmd.AddCommand(command)
}

func runVersionCommand(md *cobra.Command, args []string) {
	fmt.Println("udp-director Version: v0.1.24")
}
