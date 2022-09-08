package common

import (
	"bufio"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel
		log.Infof("[CLIENT %v] Got signal.", c.config.ID)
		c.CloseClientSocket()
	}()

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
