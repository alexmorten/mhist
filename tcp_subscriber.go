package mhist

import (
	"bufio"
	"bytes"
	"fmt"
)

//TCPSubscriber is a TCPClient that can receive messages
type TCPSubscriber struct {
	TCPClient
	newMessageChan chan []byte
}

//NewTCPSubscriber initializes a new client
func NewTCPSubscriber(address string, filterDefinition FilterDefinition, channel chan []byte) *TCPSubscriber {
	return &TCPSubscriber{
		TCPClient: TCPClient{
			Address:             address,
			buffer:              &bytes.Buffer{},
			subscriptionMessage: &SubscriptionMessage{FilterDefinition: filterDefinition},
		},
		newMessageChan: channel,
	}
}

//Read incoming messages
func (s *TCPSubscriber) Read() error {
	s.Lock()
	defer s.Unlock()
	reader := bufio.NewReader(s.conn)
	for {
		byteSlice, err := reader.ReadSlice('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		s.newMessageChan <- byteSlice
	}
}
