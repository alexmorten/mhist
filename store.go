package mhist

import (
	"sync"
)

//Store is responsible for handling Storage of different kinds of measurements
type Store struct {
	seriesMap *sync.Map
	sync.Mutex
	maxSize     int
	subscribers SubscriberSlice
	diskStore   *DiskStore
}

//NewStore ..
func NewStore(maxSize int) *Store {
	s := &Store{
		seriesMap: &sync.Map{},
		maxSize:   maxSize,
	}
	return s
}

//SetDiskStore on store
func (s *Store) SetDiskStore(ds *DiskStore) {
	s.diskStore = ds
}

//AddSubscriber to Store
func (s *Store) AddSubscriber(sub Subscriber) {
	s.subscribers = append(s.subscribers, sub)
}

//GetSeries thread safely
func (s *Store) GetSeries(name string, measurementType MeasurementType) *Series {
	series, ok := s.seriesMap.Load(name)
	if ok && series != nil {
		return series.(*Series)
	}

	s.Lock()
	defer s.Unlock()

	//Make sure we haven't added a series by chance yet
	series, ok = s.seriesMap.Load(name)
	if ok && series != nil {
		return series.(*Series)
	}
	createdSeries := NewSeries(measurementType)
	s.seriesMap.Store(name, createdSeries)
	return createdSeries
}

//Add named measurement to correct Series
func (s *Store) Add(name string, m Measurement) {
	s.subscribers.NotifyAll(name, m)

	s.GetSeries(name, m.Type()).Add(m)
}

//GetAllMeasurementsInTimeRange for all series
//TODO: change interface of diskStore to only read necessary parts, for now read all if any in memroy series is incomplete
func (s *Store) GetAllMeasurementsInTimeRange(start, end int64) map[string][]Measurement {
	m := map[string][]Measurement{}
	anyIncomplete := false

	s.forEachSeries(func(name string, series *Series) {
		measurements, incomplete := series.GetMeasurementsInTimeRange(start, end)
		m[name] = measurements
		if incomplete {
			anyIncomplete = true
		}
	})

	if s.diskStore != nil {
		allNames := s.diskStore.GetAllStoredNames()
		anyNameNotIncluded := false
		for _, name := range allNames {
			if len(m[name]) == 0 {
				anyNameNotIncluded = true
				break
			}
		}

		if anyIncomplete || anyNameNotIncluded {
			return s.diskStore.GetAllMeasurementsInTimeRange(start, end)
		}
	}

	return m
}

//Shutdown all contained series
//assumes that we don't get any messages anymore and thus don't create new Series while we do this
func (s *Store) Shutdown() {
	s.forEachSeries(func(name string, series *Series) {
		series.Shutdown()
	})
}

//Size of all carried Series'
func (s *Store) Size() int {
	size := 0

	s.forEachSeries(func(_ string, series *Series) {
		size += series.Size()
	})
	return size
}

//IsOverMaxSize we shoud start throwing things into the GC
func (s *Store) IsOverMaxSize() bool {
	return s.Size() > s.maxSize
}

//IsOverSoftLimit leaves memory room for recycling
func (s *Store) IsOverSoftLimit() bool {
	return s.Size() > int(float64(s.maxSize)*0.8)
}

//ShrinkStore by 10% and return measurements for recycling
func (s *Store) ShrinkStore() MeasurementSlices {
	slices := MeasurementSlices{}

	oldestSeries := s.findOldestSeries()
	biggestSeries := s.findBiggestSeries()
	timeRange := biggestSeries.LatestTs() - biggestSeries.OldestTs()
	cutoffPoint := oldestSeries.OldestTs() + int64(float64(timeRange)*0.1)

	s.forEachSeries(func(_ string, series *Series) {
		slices[series.Type()] = append(slices[series.Type()], series.CutoffBelow(cutoffPoint)...)
	})

	return slices
}

func (s *Store) findOldestSeries() (series *Series) {
	timestamp := int64(0)

	s.forEachSeries(func(_ string, currentSeries *Series) {
		if timestamp == 0 || timestamp > currentSeries.OldestTs() {
			timestamp = currentSeries.OldestTs()
			series = currentSeries
		}
	})
	return
}

func (s *Store) findBiggestSeries() (series *Series) {
	size := 0

	s.forEachSeries(func(_ string, currentSeries *Series) {
		if size == 0 || size < currentSeries.Size() {
			size = currentSeries.Size()
			series = currentSeries
		}
	})
	return
}

func (s *Store) forEachSeries(f func(name string, series *Series)) {
	s.seriesMap.Range(func(key, value interface{}) bool {
		name := key.(string)
		series := value.(*Series)
		if series != nil {
			f(name, series)
		}
		return true
	})
}
