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

const numericalSize = int(unsafe.Sizeof(Numerical{}))

//Size of a siggle Measurement
func (n *Numerical) Size() int {
	return numericalSize
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

//Categorical represents a single categorical meassured value in time
type Categorical struct {
	Ts    int64
	Value string
}

const categoricalSize = int(unsafe.Sizeof(Categorical{}))

//Size of a siggle Measurement
func (c *Categorical) Size() int {
	return categoricalSize + len(c.Value)
}

//Reset resets the Measurement to its zero value
func (c *Categorical) Reset() {
	c.Ts = 0
	c.Value = ""
}

//Type of Measurement
func (c *Categorical) Type() MeasurementType {
	return MeasurementCategorical
}

//Timestamp of Measurement
func (c *Categorical) Timestamp() int64 {
	return c.Ts
}

//MeasurementType enum of different types of measurements
type MeasurementType int

const (
	//MeasurementNumerical for measurements that are numerical and interpolateable
	MeasurementNumerical MeasurementType = iota

	//MeasurementCategorical for measurements that are non numerical and not interpolateable (TBD)
	MeasurementCategorical
)
