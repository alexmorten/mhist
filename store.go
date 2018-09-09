package mhist

import (
	"sync"
)

//Store is responsible for handling Storage of different kinds of measurements
type Store struct {
	seriesMap *sync.Map
	sync.Mutex
}

//NewStore ..
func NewStore() *Store {
	s := &Store{
		seriesMap: &sync.Map{},
	}
	return s
}

//GetSeries thread safely
func (s *Store) GetSeries(name string) *Series {
	series, ok := s.seriesMap.Load(name)
	if ok {
		return series.(*Series)
	}
	s.Lock()
	defer s.Unlock()

	//Make sure we haven't added a series by chance yet
	series, ok = s.seriesMap.Load(name)
	if ok {
		return series.(*Series)
	}
	createdSeries := NewSeries()
	s.seriesMap.Store(name, createdSeries)
	return createdSeries
}

//Add named measurement to correct Series
func (s *Store) Add(name string, m *Measurement) {
	s.GetSeries(name).Add(m)
}

//GetAllMeasurementsInTimeRange for all series
func (s *Store) GetAllMeasurementsInTimeRange(start, end int64) map[string][]Measurement {
	m := map[string][]Measurement{}

	s.forEachSeries(func(name string, series *Series) {
		m[name] = series.GetMeasurementsInTimeRange(start, end)
	})

	return m
}

//Shutdown all contained series
//assumes that we don't get any messages anymore and thus don't create new Series while we do this
func (s *Store) Shutdown() {
	s.forEachSeries(func(name string, series *Series) {
		series.Shutdown()
	})
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
