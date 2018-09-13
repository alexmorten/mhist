package mhist

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

//Series represents a series of measurements over time
//assumes measurements are taken in order
type Series struct {
	measurements    []Measurement
	addChan         chan Measurement
	cutoffChan      chan *cutoffMessage
	stopChan        chan struct{}
	size            int
	measurementType MeasurementType
	rwLock          sync.RWMutex
}

type cutoffMessage struct {
	lowestTs   int64
	returnChan chan []Measurement
}

//NewSeries constructs a new series and starts the listening goroutine
func NewSeries(measurementType MeasurementType) *Series {
	s := &Series{
		measurements:    []Measurement{},
		addChan:         make(chan Measurement),
		stopChan:        make(chan struct{}),
		cutoffChan:      make(chan *cutoffMessage),
		measurementType: measurementType,
	}

	go s.Listen()
	return s
}

//Add m to series
func (s *Series) Add(m Measurement) {
	s.addChan <- m
}

//CutoffBelow a timestamp and return thrown away measurements
func (s *Series) CutoffBelow(lowestTs int64) []Measurement {
	returnChan := make(chan []Measurement)
	s.cutoffChan <- &cutoffMessage{
		lowestTs:   lowestTs,
		returnChan: returnChan,
	}
	returnedSlice := <-returnChan
	return returnedSlice
}

//Shutdown series goroutine
func (s *Series) Shutdown() {
	s.stopChan <- struct{}{}
}

//GetMeasurementsInTimeRange returns the measurements in approx. the given timerange
////assumes equally distributed measurements over time
func (s *Series) GetMeasurementsInTimeRange(start int64, end int64) []Measurement {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
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
		measurements[i] = getCopiedValueInterface(s.measurements[i+startIndex])
	}
	return measurements
}

//Listen for new measurements
func (s *Series) Listen() {
loop:
	for {
		select {
		case <-s.stopChan:
			break loop
		case message := <-s.cutoffChan:
			s.handleCutoff(message)
		case measurement := <-s.addChan:
			s.handleAdd(measurement)
		}
	}
}

//Size of all measurements contained in the Series
func (s *Series) Size() int {
	return s.size
}

//Type of contained measurements
func (s *Series) Type() MeasurementType {
	return s.measurementType
}

func (s *Series) handleCutoff(message *cutoffMessage) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	if message.lowestTs <= s.OldestTs() {
		message.returnChan <- []Measurement{}
		return
	}

	index := 0
	removedBytes := 0
	for _, m := range s.measurements {
		if m.Timestamp() > message.lowestTs {
			break
		}
		removedBytes += m.Size()
		index++
	}

	cutoffSlices := s.measurements[:index]
	remainingSlices := s.measurements[index:]
	s.measurements = remainingSlices
	s.size -= removedBytes
	message.returnChan <- cutoffSlices
}

func (s *Series) handleAdd(m Measurement) {
	s.size += m.Size()
	s.measurements = append(s.measurements, m)
}

func (s *Series) calcIndexAbove(ts int64) (int, error) {
	if ts <= s.OldestTs() {
		return 0, nil
	}
	//shouldn't happen
	if ts > s.LatestTs() {
		return 0, errors.New("given ts is above the latest measured timestamp")
	}

	//assumes equally distributed measurements over time, no need for perfectly accurate results yet
	timeRange := s.LatestTs() - s.OldestTs()
	posInRange := ts - s.OldestTs()
	index := float64(posInRange) / float64(timeRange) * float64(len(s.measurements)-1)
	return int(math.Ceil(index)), nil
}

func (s *Series) calcIndexBelow(ts int64) (int, error) {
	//shouldn't happen
	if ts < s.OldestTs() {
		return 0, errors.New("given ts is below the oldest measured timestamp")
	}
	if ts >= s.LatestTs() {
		return len(s.measurements) - 1, nil
	}

	//assumes equally distributed measurements over time, no need for perfectly accurate results yet
	timeRange := s.LatestTs() - s.OldestTs()
	posInRange := ts - s.OldestTs()
	index := float64(posInRange) / float64(timeRange) * float64(len(s.measurements)-1)
	return int(index), nil
}

//LatestTs in series
func (s *Series) LatestTs() int64 {
	if len(s.measurements) == 0 {
		return 0
	}

	return s.measurements[len(s.measurements)-1].Timestamp()
}

//OldestTs in series
func (s *Series) OldestTs() int64 {
	if len(s.measurements) == 0 {
		return 0
	}

	return s.measurements[0].Timestamp()
}

func getCopiedValueInterface(measurement Measurement) Measurement {
	switch measurement.(type) {
	case *Numerical:
		value := *measurement.(*Numerical)
		return &value
	}
	return nil
}
