package mhist

import (
	"sync"

	"github.com/alexmorten/mhist/models"
)

type notifyMessage struct {
	name        string
	measurement models.Measurement
}

type grpcSubscriber struct {
	notifyChan chan notifyMessage
}

func (s *grpcSubscriber) Notify(name string, measurement models.Measurement) {
	s.notifyChan <- notifyMessage{
		name:        name,
		measurement: measurement,
	}
}

func (s *grpcSubscriber) drain() {
	go func() {
		for range s.notifyChan {
		}
	}()
}

type grpcSubscribers struct {
	list []*grpcSubscriber
	*sync.RWMutex
}

func newGrpcSubscribers() *grpcSubscribers {
	return &grpcSubscribers{
		RWMutex: &sync.RWMutex{},
	}
}

func (subs *grpcSubscribers) newSubscriber() *grpcSubscriber {
	s := &grpcSubscriber{
		notifyChan: make(chan notifyMessage),
	}

	subs.Lock()
	subs.list = append(subs.list, s)
	defer subs.Unlock()

	return s
}

func (subs *grpcSubscribers) removeSubscriber(subscriberToDelete *grpcSubscriber) {
	subscriberToDelete.drain()
	subs.Lock()
	defer subs.Unlock()
	if len(subs.list) == 0 {
		return
	}

	remainingList := make([]*grpcSubscriber, 0, len(subs.list)-1)
	for _, s := range subs.list {
		if s != subscriberToDelete {
			remainingList = append(remainingList, s)
		}
	}

	subs.list = remainingList

	close(subscriberToDelete.notifyChan)
}

func (subs *grpcSubscribers) forEach(f func(*grpcSubscriber)) {
	subs.RLock()
	defer subs.RUnlock()

	for _, s := range subs.list {
		f(s)
	}
}
