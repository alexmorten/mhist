package models

//SubscriptionMessage is the message the client sends to the server
type SubscriptionMessage struct {
	Publisher        bool             `json:"publisher"`
	FilterDefinition FilterDefinition `json:"filter"`
}
