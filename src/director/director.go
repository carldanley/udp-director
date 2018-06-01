package director

import (
	"fmt"
	"net"
	"time"

	"github.com/carldanley/udp-director/src/utils"
)

type Director struct {
	ActiveConnections             map[string]*Connection
	UDPServer                     *net.UDPConn
	ErrorChannel                  chan error
	InactiveConnectionTimeSeconds float64
	Destinations                  []utils.ParsedDestination
}

func createUDPServer(sourceHost string, sourcePort int) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", sourceHost, sourcePort))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateNewDirector(sourceHost string, sourcePort int, destinations []utils.ParsedDestination, inactiveConnectionTimeSeconds float64) (*Director, error) {
	udpServer, err := createUDPServer(sourceHost, sourcePort)
	if err != nil {
		return nil, err
	}

	director := Director{
		ActiveConnections:             map[string]*Connection{},
		UDPServer:                     udpServer,
		ErrorChannel:                  make(chan error),
		InactiveConnectionTimeSeconds: inactiveConnectionTimeSeconds,
		Destinations:                  destinations,
	}

	for _, destination := range destinations {
		fmt.Printf("Forwarding %s:%d -> %s:%d\n", sourceHost, sourcePort, destination.IP, destination.Port)
	}
	fmt.Println("")

	return &director, nil
}

func (d *Director) Stop() {
	if d.UDPServer != nil {
		d.UDPServer.Close()
		d.UDPServer = nil
	}

	for _, connection := range d.ActiveConnections {
		connection.Stop()
	}

	d.ActiveConnections = nil

	if d.ErrorChannel != nil {
		close(d.ErrorChannel)
		d.ErrorChannel = nil
	}
}

func (d *Director) Listen() {
	buffer := make([]byte, 1024)

	go d.expireInactiveConnections()

	for {
		bytesSent, remoteAddress, err := d.UDPServer.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		remoteAddressStr := remoteAddress.String()

		if connection, found := d.ActiveConnections[remoteAddressStr]; found {
			connection.onIngressReceived(buffer[0:bytesSent], d)
		} else {
			fmt.Printf("[ACTIVE]: %s...\n", remoteAddressStr)

			connection, err := createNewConnection(remoteAddressStr, d)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			d.ActiveConnections[remoteAddressStr] = connection
			connection.onIngressReceived(buffer[0:bytesSent], d)
		}
	}
}

func (d *Director) expireInactiveConnections() {
	activeConnections := map[string]*Connection{}

	for remoteAddress, connection := range d.ActiveConnections {
		elapsedTime := time.Since(connection.LastActivity)

		if elapsedTime.Seconds() < d.InactiveConnectionTimeSeconds {
			activeConnections[remoteAddress] = connection
		} else {
			fmt.Printf("[INACTIVE]: %s\n", remoteAddress)
			connection.Stop()
		}
	}

	d.ActiveConnections = activeConnections

	time.Sleep(time.Second)
	d.expireInactiveConnections()
}
