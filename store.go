package mhist

import (
	"github.com/alexmorten/mhist/models"
)

//Store is responsible for handling Storage of different kinds of measurements
type Store struct {
	subscribers SubscriberSlice
	diskStore   *DiskStore
}

//NewStore from diskstore, that handles subscribers
func NewStore(diskStore *DiskStore) *Store {
	store := &Store{
		diskStore: diskStore,
	}
	store.AddSubscriber(diskStore)
	return store
}

//AddSubscriber to Store
func (s *Store) AddSubscriber(sub Subscriber) {
	s.subscribers = append(s.subscribers, sub)
}

//Add named measurement
func (s *Store) Add(name string, m models.Measurement) {
	s.subscribers.NotifyAll(name, m)
}

//GetMeasurementsInTimeRange from disk store
func (s *Store) GetMeasurementsInTimeRange(start, end int64, filterDefinition models.FilterDefinition) map[string][]models.Measurement {
	return s.diskStore.GetMeasurementsInTimeRange(start, end, filterDefinition)
}

//GetStoredMetaInfo from Diskstore
func (s *Store) GetStoredMetaInfo() []MeasurementTypeInfo {
	return s.diskStore.GetAllStoredInfos()
}

//Shutdown diskStore
func (s *Store) Shutdown() {
	s.diskStore.Shutdown()
}
