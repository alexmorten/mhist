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
	pools := &Pools{
		Store: store,
	}
	pools.pools = map[MeasurementType]*sync.Pool{
		MeasurementNumerical: &sync.Pool{
			New: func() interface{} {
				slices, ok := grabSlicesFromStore(store)
				if ok {
					numericalSlice := slices[MeasurementNumerical]
					if len(numericalSlice) > 0 {
						measurement := numericalSlice[0]
						rest := numericalSlice[1:]
						slices[MeasurementNumerical] = rest
						pools.fill(slices)
						return measurement
					}
				}
				return &Measurement{}
			},
		},
	}
	return pools
}

//GetNumericalMeasurement out of the correct pool
func (pools *Pools) GetNumericalMeasurement() *Measurement {
	return pools.pools[MeasurementNumerical].Get().(*Measurement)
}

func grabSlicesFromStore(store *Store) (slices MeasurementSlices, ok bool) {
	if store.IsOverMaxSize() {
		return store.ShrinkStore(), true
	}
	return nil, false
}

func (pools *Pools) fill(slices MeasurementSlices) {
	for key, slice := range slices {
		for _, measurement := range slice {
			pools.pools[key].Put(measurement)
		}
	}
}
