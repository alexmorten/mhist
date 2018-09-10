package mhist

import (
	"unsafe"
)

//Measurement represents a single meassured value in time
type Measurement struct {
	Ts    int64
	Value float64
}

const measurementSize = int(unsafe.Sizeof(Measurement{}))

//Size of a siggle Measurement
func (m *Measurement) Size() int {
	return measurementSize
}

//Reset resets the Measurement to its zero value
func (m *Measurement) Reset() {
	m.Ts = 0
	m.Value = 0
}

//Type of Measurement
func (m *Measurement) Type() MeasurementType {
	return MeasurementNumerical
}

//MeasurementType enum of different types of measurements
type MeasurementType int

const (
	//MeasurementNumerical for measurements that are numerical and interpolateable
	MeasurementNumerical MeasurementType = iota

	//MeasurementCategorical for measurements that are non numerical and not interpolateable (TBD)
	// MeasurementCategorical
)
