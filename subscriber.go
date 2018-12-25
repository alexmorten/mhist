package mhist

import "github.com/alexmorten/mhist/models"

//Subscriber is something that can be notified with measurements
type Subscriber interface {
	Notify(name string, measurement models.Measurement)
}

//SubscriberSlice is a collection of Subscribers, with helper methods
type SubscriberSlice []Subscriber

//NotifyAll Subscribers in Slice
func (s SubscriberSlice) NotifyAll(name string, measurement models.Measurement) {
	for _, subscriber := range s {
		subscriber.Notify(name, measurement)
	}
}
