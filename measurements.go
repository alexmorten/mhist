package mhist

import (
	"unsafe"
)

//Measurement interface
type Measurement interface {
	Type() MeasurementType
	Timestamp() int64
	Size() int
	Reset()
}

//Numerical represents a single meassured value in time
type Numerical struct {
	Ts    int64
	Value float64
}

const measurementSize = int(unsafe.Sizeof(Numerical{}))

//Size of a siggle Measurement
func (n *Numerical) Size() int {
	return measurementSize
}

//Reset resets the Measurement to its zero value
func (n *Numerical) Reset() {
	n.Ts = 0
	n.Value = 0
}

//Type of Measurement
func (n *Numerical) Type() MeasurementType {
	return MeasurementNumerical
}

//Timestamp of Measurement
func (n *Numerical) Timestamp() int64 {
	return n.Ts
}

//MeasurementType enum of different types of measurements
type MeasurementType int

const (
	//MeasurementNumerical for measurements that are numerical and interpolateable
	MeasurementNumerical MeasurementType = iota

	//MeasurementCategorical for measurements that are non numerical and not interpolateable (TBD)
	// MeasurementCategorical
)
