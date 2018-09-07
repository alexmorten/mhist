package mhist

import (
	"errors"
	"fmt"
	"math"
)

//Series represents a series of measurements over time
//assumes measurements are taken in order
type Series struct {
	measurements []*Measurement
	addChan      chan *Measurement
	stopChan     chan struct{}
}

//NewSeries constructs a new series and starts the listening goroutine
func NewSeries() *Series {
	s := &Series{
		measurements: []*Measurement{},
		addChan:      make(chan *Measurement),
		stopChan:     make(chan struct{}),
	}
	go s.Listen()
	return s
}

//Add m to series
func (s *Series) Add(m *Measurement) {
	s.addChan <- m
}

//Shutdown series goroutine
func (s *Series) Shutdown() {
	s.stopChan <- struct{}{}
}

//GetMeasurementsInTimeRange returns the measurements in approx. the given timerange
////assumes equally distributed measurements over time
func (s *Series) GetMeasurementsInTimeRange(start int64, end int64) []Measurement {
	startIndex, err := s.calcIndexAbove(start)
	if err != nil {
		fmt.Println(err)
		return []Measurement{}
	}
	endIndex, err := s.calcIndexBelow(end)
	if err != nil {
		fmt.Println(err)
		return []Measurement{}
	}
	length := endIndex - startIndex + 1
	measurements := make([]Measurement, length)
	for i := 0; i < length; i++ {
		measurements[i] = *s.measurements[i+startIndex]
	}
	return measurements
}

//Listen for new measurements
func (s *Series) Listen() {
loop:
	for {
		select {
		case m := <-s.addChan:
			s.handleAdd(m)
		case <-s.stopChan:
			break loop
		}
	}
}

func (s *Series) handleAdd(m *Measurement) {
	s.measurements = append(s.measurements, m)
}

func (s *Series) calcIndexAbove(ts int64) (int, error) {
	if ts <= s.oldestTs() {
		return 0, nil
	}
	//shouldn't happen
	if ts > s.latestTs() {
		return 0, errors.New("given ts is above the latest measured timestamp")
	}

	//assumes equally distributed measurements over time, no need for perfectly accurate results yet
	timeRange := s.latestTs() - s.oldestTs()
	posInRange := ts - s.oldestTs()
	index := float64(posInRange) / float64(timeRange) * float64(len(s.measurements)-1)
	return int(math.Ceil(index)), nil
}

func (s *Series) calcIndexBelow(ts int64) (int, error) {
	//shouldn't happen
	if ts < s.oldestTs() {
		return 0, errors.New("given ts is below the oldest measured timestamp")
	}
	if ts >= s.latestTs() {
		return len(s.measurements) - 1, nil
	}

	//assumes equally distributed measurements over time, no need for perfectly accurate results yet
	timeRange := s.latestTs() - s.oldestTs() // 40
	posInRange := ts - s.oldestTs()          // 1035 - 1000 = 35
	index := float64(posInRange) / float64(timeRange) * float64(len(s.measurements)-1)
	return int(index), nil
}

func (s *Series) latestTs() int64 {
	if len(s.measurements) == 0 {
		return 0
	}

	return s.measurements[len(s.measurements)-1].Ts
}

func (s *Series) oldestTs() int64 {
	if len(s.measurements) == 0 {
		return 0
	}

	return s.measurements[0].Ts
}
