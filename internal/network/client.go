package network

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	network string
	address string
}

func CreateClient(network string, host string, port int) *Client {
	return &Client{
		network: network,
		address: fmt.Sprintf("%s:%d", host, port),
	}
}

func (c *Client) Call(message string) string {
	connection, err := net.Dial(c.network, c.address)
	if err != nil {
		log.Fatalf("Cannot connect to server '%s'", c.address)
	}

	defer connection.Close()

	data := []byte(message)
	_, err = connection.Write(data)
	if err != nil {
		log.Fatalf("cannot send data '%s' to server", data)
	}

	response := make([]byte, 1024)
	end, err := connection.Read(response)
	if err != nil {
		log.Fatalf("Cannot recieve data from server, error: %s", err)
	}

	return string(response[:end])
}
