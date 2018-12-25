package models

//Message represents events sent to and from the server
type Message struct {
	Name      string      `json:"name"`
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}

//Reset message to zero value
func (m *Message) Reset() {
	m.Name = ""
	m.Timestamp = 0
	m.Value = nil
}
