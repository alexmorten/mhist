package mhist

import (
	"encoding/json"
	"fmt"
)

//Replication is a wrapper for tcp.Client, that implements the subscriber interface
type Replication struct {
	client *TCPClient
}

//NewReplication creates the underlying tcp.Client correctly
func NewReplication(address string) *Replication {
	return &Replication{
		client: NewReplicatorClient(address),
	}
}

//Notify replication about new measurement
func (r *Replication) Notify(name string, measurement Measurement) {
	message := &message{
		Name:  name,
		Value: measurement.ValueInterface(),
	}
	byteSlice, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	go r.client.Write(byteSlice)
}
