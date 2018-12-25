package models

import (
	"sync"
)

//Pools holds the pools for the different measurement types (only one for now)
type Pools struct {
	messagePool      *sync.Pool
	measurementPools map[MeasurementType]*sync.Pool
}

//MeasurementSlices for moving types of Measurements around
type MeasurementSlices map[MeasurementType][]Measurement

//NewPools returns the constructed pool handler
func NewPools() *Pools {
	pools := &Pools{
		messagePool: &sync.Pool{
			New: func() interface{} {
				return &Message{}
			},
		},
	}
	pools.measurementPools = map[MeasurementType]*sync.Pool{
		MeasurementNumerical: &sync.Pool{
			New: func() interface{} {
				return &Numerical{}
			},
		},
		MeasurementCategorical: &sync.Pool{
			New: func() interface{} {
				return &Categorical{}
			},
		},
	}
	return pools
}

//GetNumericalMeasurement out of the correct pool
func (pools *Pools) GetNumericalMeasurement() *Numerical {
	return pools.measurementPools[MeasurementNumerical].Get().(*Numerical)
}

//GetCategoricalMeasurement out of the correct pool
func (pools *Pools) GetCategoricalMeasurement() *Categorical {
	return pools.measurementPools[MeasurementCategorical].Get().(*Categorical)
}

//PutMeasurement out of the correct pool
func (pools *Pools) PutMeasurement(m Measurement) {
	pools.measurementPools[m.Type()].Put(m)
}

//GetMessage from MessagePool
func (pools *Pools) GetMessage() *Message {
	return pools.messagePool.Get().(*Message)
}

//PutMessage into MessagePool
func (pools *Pools) PutMessage(m *Message) {
	pools.messagePool.Put(m)
}
