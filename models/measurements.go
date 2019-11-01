package models

//Measurement interface
type Measurement interface {
	Type() MeasurementType
	Timestamp() int64
}

//Numerical represents a single meassured value in time
type Numerical struct {
	Ts    int64
	Value float64
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

//Type of Measurement
func (c *Categorical) Type() MeasurementType {
	return MeasurementCategorical
}

//Timestamp of Measurement
func (c *Categorical) Timestamp() int64 {
	return c.Ts
}

// Raw represents any measurements that are neither categorical nor a single number, by simply treating them as bytes
type Raw struct {
	Ts    int64
	Value []byte
}

//Type of Measurement
func (r *Raw) Type() MeasurementType {
	return MeasurementRaw
}

//Timestamp of Measurement
func (r *Raw) Timestamp() int64 {
	return r.Ts
}

//MeasurementType enum of different types of measurements
type MeasurementType int

const (
	//MeasurementNumerical for measurements that are numerical and interpolateable
	MeasurementNumerical MeasurementType = iota + 1 //start at 1 so we can differentiate from zero value

	//MeasurementCategorical for measurements that are non numerical and not interpolateable
	MeasurementCategorical

	// MeasurementRaw for measurements that can only be represented as raw bytes
	MeasurementRaw
)
