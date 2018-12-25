package tcp

import (
	"encoding/json"
	"fmt"

	"github.com/alexmorten/mhist/models"
)

//Replication is a wrapper for tcp.Client, that implements the subscriber interface
type Replication struct {
	client *Client
	pools  *models.Pools
}

//NewReplication creates the underlying tcp.Client correctly
func NewReplication(address string, pools *models.Pools) *Replication {
	return &Replication{
		client: NewReplicatorClient(address),
		pools:  pools,
	}
}

//Notify replication about new measurement
func (r *Replication) Notify(name string, measurement models.Measurement) {
	message := r.pools.GetMessage()
	defer r.pools.PutMessage(message)

	message.Reset()
	message.Name = name
	message.Value = measurement.ValueInterface()
	message.Timestamp = measurement.Timestamp()

	byteSlice, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	go r.client.Write(byteSlice)

}
