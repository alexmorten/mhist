package mhist

import (
	"sync"
)

//Pools holds the pools for the different measurement types (only one for now)
type Pools struct {
	Store *Store
	pools map[MeasurementType]*sync.Pool
}

//MeasurementSlices for moving types of Measurements around
type MeasurementSlices map[MeasurementType][]*Measurement

//NewPools returns the constructed pool handler
func NewPools(store *Store) *Pools {
	return &Pools{
		Store: store,
		pools: map[MeasurementType]*sync.Pool{
			MeasurementNumerical: &sync.Pool{
				New: func() interface{} {
					return struct{}{}
				},
			},
		},
	}
}

//GrabSlicesFromStore ...
func (pools *Pools) GrabSlicesFromStore() MeasurementSlices {
	return MeasurementSlices{}
}

//Fill pools with slices
func (pools *Pools) Fill(slices MeasurementSlices) {
	for key, slice := range slices {
		for _, measurement := range slice {
			pools.pools[key].Put(measurement)
		}
	}
}
