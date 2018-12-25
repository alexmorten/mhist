package tcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/alexmorten/mhist/models"
)

//Client for tcp connections. Automatically retries establishing connections
type Client struct {
	Address             string
	subscriptionMessage *models.SubscriptionMessage
	buffer              *bytes.Buffer
	conn                net.Conn
	sync.RWMutex
}

//NewClient initializes a new client
func NewClient(address string) *Client {
	return &Client{
		Address:             address,
		buffer:              &bytes.Buffer{},
		subscriptionMessage: &models.SubscriptionMessage{Publisher: true},
	}
}

//NewReplicatorClient sets the subscriptionMessage correctly for a replication connection
func NewReplicatorClient(address string) *Client {
	client := NewClient(address)
	client.subscriptionMessage.Replication = true
	return client
}

//Connect to described address
func (c *Client) Connect() {
	c.Lock()
	defer c.Unlock()
	for {
		conn, err := net.Dial("tcp", c.Address)
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		message, err := json.Marshal(c.subscriptionMessage)
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(append(message, '\n'))
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		c.conn = conn
		if c.buffer.Len() > 0 {
			c.writeBufferToConnection()
		}
		break
	}
}

func (c *Client) Write(byteSlice []byte) {
	c.Lock()
	c.buffer.Write(append(byteSlice, '\n'))
	err := c.writeBufferToConnection()
	if err != nil && c.conn == nil {
		c.Unlock()
		go c.Connect()
	} else {
		c.Unlock()
	}
}

func (c *Client) writeBufferToConnection() error {
	if c.conn == nil {
		return errors.New("connection not set")
	}
	_, err := c.conn.Write(c.buffer.Bytes())
	if err == nil {
		c.buffer.Reset()
	} else {
		c.conn = nil
	}
	return err
}
