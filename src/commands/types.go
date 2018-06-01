package commands

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "udp-director",
	Short: "udp-director is a tool that can handle redirecting UDP traffic to any number of servers",
}

var sourceHost string
var sourcePort int
var destinations string
var inactiveConnectionTimeSeconds float64
