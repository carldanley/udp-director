package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedDestination struct {
	IP   string
	Port int
}

func ParseDestinations(str string) ([]ParsedDestination, error) {
	destinations := []ParsedDestination{}

	parts := strings.Split(str, ",")

	for _, part := range parts {
		tmp := strings.Split(part, ":")

		if len(tmp) != 2 {
			return nil, fmt.Errorf("Invalid destination was specified: %s", part)
		}

		port, err := strconv.Atoi(tmp[1])
		if err != nil {
			return nil, fmt.Errorf("An invalid port was specified: %s", part)
		}

		destination := ParsedDestination{
			IP:   tmp[0],
			Port: port,
		}

		destinations = append(destinations, destination)
	}

	return destinations, nil
}
