package models

import (
	"strconv"
	"unsafe"
)

//Measurement interface
type Measurement interface {
	Type() MeasurementType
	Timestamp() int64
	ValueString() string
	ValueInterface() interface{}
	Size() int
	Reset()
	Copy() Measurement
	CopyFrom(p *Pools) Measurement
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

//ValueInterface of Measurement
func (n *Numerical) ValueInterface() interface{} {
	return n.Value
}

//ValueString of Measurement
func (n *Numerical) ValueString() string {
	return strconv.FormatFloat(n.Value, 'f', -1, 64)
}

//Copy of Measurement
func (n *Numerical) Copy() Measurement {
	value := *n
	return &value
}

//CopyFrom pools for GC efficiency
func (n *Numerical) CopyFrom(p *Pools) Measurement {
	newNumerical := p.GetNumericalMeasurement()
	newNumerical.Ts = n.Ts
	newNumerical.Value = n.Value
	return newNumerical
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

//ValueInterface of Measurement
func (c *Categorical) ValueInterface() interface{} {
	return c.Value
}

//ValueString of Measurement
func (c *Categorical) ValueString() string {
	return c.Value
}

//Copy of Measurement
func (c *Categorical) Copy() Measurement {
	value := *c
	return &value
}

//CopyFrom pools for GC efficiency
func (c *Categorical) CopyFrom(p *Pools) Measurement {
	newCategorical := p.GetCategoricalMeasurement()
	newCategorical.Ts = c.Ts
	newCategorical.Value = c.Value
	return newCategorical
}

//MeasurementType enum of different types of measurements
type MeasurementType int

const (
	//MeasurementNumerical for measurements that are numerical and interpolateable
	MeasurementNumerical MeasurementType = iota + 1 //start at 1 so we can differenciate from zero value

	//MeasurementCategorical for measurements that are non numerical and not interpolateable
	MeasurementCategorical
)
