package mhist

import (
	"encoding/json"
	"fmt"

	"github.com/codeuniversity/ppp-mhist/tcp"
)

//Replication is a wrapper for tcp.Client, that implements the subscriber interface
type Replication struct {
	client *tcp.Client
}

//NewReplication creates the underlying tcp.Client correctly
func NewReplication(address string) *Replication {
	return &Replication{
		client: tcp.NewReplicatorClient(address),
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
	r.client.Write(byteSlice)
}
