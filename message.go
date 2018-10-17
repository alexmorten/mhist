package mhist

//Message represents events sent to and from the server
type Message struct {
	Name      string      `json:"name"`
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}
