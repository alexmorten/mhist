package mhist

//Subscriber is something that can be notified with measurements
type Subscriber interface {
	Notify(name string, measurement Measurement)
}

//SubscriberSlice is a collection of Subscribers, with helper methods
type SubscriberSlice []Subscriber

//NotifyAll Subscribers in Slice
func (s SubscriberSlice) NotifyAll(name string, measurement Measurement) {
	for _, subscriber := range s {
		subscriber.Notify(name, measurement)
	}
}
