package common

import (
	"bufio"
	"net"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) CreateClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"[CLIENT %v] Could not connect to server. Error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) CloseClientSocket() error {
	c.conn.Close()
	return nil
}

func (c *Client) queryWinners(contestants *ContestantsQuery) *WinnerResponse {
	writer := bufio.NewWriter(c.conn)
	WriteContestantsMessage(contestants, writer)
	writer.Flush()

	reader := bufio.NewReader(c.conn)
	return ReadWinnerMessage(reader)
}
