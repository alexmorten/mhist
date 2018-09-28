package tcp

//SubscribtionMessage is the message the client sends to the server to make sure we can use the same logic for realtime update streams and replications
//where any server in the cluster can be a listening point for realtime updates, without having endless replication messages bouncing between the servers
type SubscribtionMessage struct {
	Replication bool `json:"replication"`
	Publisher   bool `json:"publisher"`
}
