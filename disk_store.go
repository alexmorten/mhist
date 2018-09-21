package mhist

import (
	"fmt"
	"os"
	"sync"
)

var dataPath = "data"

//DiskStore handles writes and reads from disk
type DiskStore struct {
	pools    *Pools
	blockMap *sync.Map
	sync.Mutex
}

//NewDiskStore initializes the disk store
func NewDiskStore(pools *Pools) *DiskStore {
	os.MkdirAll(dataPath, os.ModePerm)
	return &DiskStore{
		pools:    pools,
		blockMap: &sync.Map{},
	}
}

//Notify DiskStore about new Measurement
func (s *DiskStore) Notify(name string, m Measurement) {
	ownMeasurement := m.CopyFrom(s.pools)

	fmt.Println(ownMeasurement)
	block, err := s.GetBlock(name, ownMeasurement.Type())
	if err != nil {
		fmt.Println(err)
		return
	}
	block.Add(ownMeasurement)
	s.pools.PutMeasurement(ownMeasurement)
}

//GetBlock thread safely
func (s *DiskStore) GetBlock(name string, measurementType MeasurementType) (*DiskBlock, error) {
	block, ok := s.blockMap.Load(name)
	if ok && block != nil {
		return block.(*DiskBlock), nil
	}

	s.Lock()
	defer s.Unlock()

	//Make sure we haven't added a block by chance yet
	block, ok = s.blockMap.Load(name)
	if ok && block != nil {
		return block.(*DiskBlock), nil
	}
	createdBlock, err := NewDiskBlock(measurementType, name)
	if err != nil {
		return nil, err
	}
	s.blockMap.Store(name, createdBlock)
	return createdBlock, nil
}
