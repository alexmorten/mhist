package mhist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

//TCPClient for tcp connections. Automatically retries establishing connections
type TCPClient struct {
	Address             string
	subscriptionMessage *SubscriptionMessage
	buffer              *bytes.Buffer
	conn                net.Conn
	sync.RWMutex
}

//NewTCPClient initializes a new client
func NewTCPClient(address string) *TCPClient {
	return &TCPClient{
		Address:             address,
		buffer:              &bytes.Buffer{},
		subscriptionMessage: &SubscriptionMessage{Publisher: true},
	}
}

//NewReplicatorClient sets the subscriptionMessage correctly for a replication connection
func NewReplicatorClient(address string) *TCPClient {
	client := NewTCPClient(address)
	client.subscriptionMessage.Replication = true
	return client
}

//Connect to described address
func (c *TCPClient) Connect() {
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

func (c *TCPClient) Write(byteSlice []byte) {
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

func (c *TCPClient) writeBufferToConnection() error {
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
