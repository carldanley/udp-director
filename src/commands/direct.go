package commands

import (
	"fmt"
	"os"

	"github.com/carldanley/udp-director/src/director"
	"github.com/carldanley/udp-director/src/utils"
	"github.com/spf13/cobra"
)

func init() {
	var command = &cobra.Command{
		Use:   "direct",
		Short: "Directs traffic from a port to other IP/port combinations",
		Run:   runDirectCommand,
	}

	// add flags for this command
	command.Flags().StringVarP(&sourceHost, "source-host", "s", "0.0.0.0", "Specifies the host address to listen for UDP traffic on")
	command.Flags().IntVarP(&sourcePort, "source-port", "p", 1337, "Specifies a port to listen for UDP traffic on")
	command.Flags().StringVarP(&destinations, "destinations", "d", "", "Specifies a comma-seperated list of IP:Port combinations to direct traffic to")
	command.Flags().Float64VarP(&inactiveConnectionTimeSeconds, "inactive-connection-time-seconds", "i", 10, "Specifies the number of seconds required for activity until an incoming connection is expired")

	// register this command
	RootCmd.AddCommand(command)
}

func runDirectCommand(md *cobra.Command, args []string) {
	destinations, err := utils.ParseDestinations(destinations)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(destinations) == 0 {
		fmt.Println("Please specify at least one destination to direct UDP traffic to")
		os.Exit(1)
	}

	director, err := director.CreateNewDirector(sourceHost, sourcePort, destinations, inactiveConnectionTimeSeconds)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer director.Stop()

	go director.Listen()

	for err := range director.ErrorChannel {
		fmt.Println(err.Error())
	}
}
