package tcp

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/alexmorten/mhist/models"
)

//Subscriber is a TCPClient that can receive messages
type Subscriber struct {
	Client
	newMessageChan chan []byte
}

//NewTCPSubscriber initializes a new client
func NewTCPSubscriber(address string, filterDefinition models.FilterDefinition, channel chan []byte) *Subscriber {
	return &Subscriber{
		Client: Client{
			Address:             address,
			buffer:              &bytes.Buffer{},
			subscriptionMessage: &models.SubscriptionMessage{FilterDefinition: filterDefinition},
		},
		newMessageChan: channel,
	}
}

//Read incoming messages
func (s *Subscriber) Read() error {
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
