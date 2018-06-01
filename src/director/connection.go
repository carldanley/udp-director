package director

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Connection struct {
	Address      string
	LastActivity time.Time

	director *Director
	relayers []*net.UDPConn
}

func createUDPClient(destination string) (*net.UDPConn, error) {
	serverAddress, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func createNewConnection(remoteAddress string, director *Director) (*Connection, error) {
	connection := Connection{
		Address:      remoteAddress,
		LastActivity: time.Now(),
		director:     director,
		relayers:     []*net.UDPConn{},
	}

	for _, destination := range director.Destinations {
		client, err := createUDPClient(fmt.Sprintf("%s:%d", destination.IP, destination.Port))
		if err != nil {
			return nil, err
		}

		go connection.ListenForEgress(client)

		connection.relayers = append(connection.relayers, client)
	}

	return &connection, nil
}

func (c *Connection) ListenForEgress(client *net.UDPConn) {
	for {
		if c.relayers == nil {
			return
		}

		buffer := make([]byte, 1024)

		bytesSent, _, err := client.ReadFromUDP(buffer)
		if err != nil {
			if c.relayers != nil {
				continue
			} else if strings.Contains(err.Error(), "recvfrom: connection refused") {
				continue
			}

			fmt.Println(err.Error())
			continue
		}

		c.onEgressReceived(buffer[0:bytesSent])
	}
}

func (c *Connection) onIngressReceived(payload []byte, director *Director) {
	c.LastActivity = time.Now()

	for _, relayer := range c.relayers {
		_, err := relayer.Write(payload)

		if err != nil {
			director.ErrorChannel <- err
		}
	}
}

func (c *Connection) onEgressReceived(payload []byte) {
	addr, err := net.ResolveUDPAddr("udp", c.Address)
	if err != nil {
		c.director.ErrorChannel <- err
		return
	}

	if _, err := c.director.UDPServer.WriteTo(payload, addr); err != nil {
		c.director.ErrorChannel <- err
		return
	}
}

func (c *Connection) Stop() {
	relayers := c.relayers
	c.relayers = nil

	for _, relayer := range relayers {
		relayer.Close()
	}
}
