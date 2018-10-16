package mhist

import (
	"bufio"
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
			subscriptionMessage: &SubscriptionMessage{FilterDefinition: filterDefinition},
		},
		newMessageChan: channel,
	}
}

//Read incoming messages
func (s *TCPSubscriber) Read() error {
	s.Lock()
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
