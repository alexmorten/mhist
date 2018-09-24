package mhist

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maxBuffer = 12 * 1024
const maxFileSize = 10 * 1024 * 1024
const maxDiskSize = 1 * 1024 * 1024 * 1024

var dataPath = "data"

//DiskStore handles buffered writes to Disk
type DiskStore struct {
	block    *Block
	meta     *DiskMeta
	pools    *Pools
	addChan  chan addMessage
	stopChan chan struct{}
}

type addMessage struct {
	name        string
	measurement Measurement
	doneChan    chan struct{}
}

//NewDiskStore initializes the DiskBlockRoutine
func NewDiskStore(pools *Pools) (*DiskStore, error) {
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	block := &DiskStore{
		meta:     InitMetaFromDisk(),
		block:    &Block{},
		addChan:  make(chan addMessage),
		stopChan: make(chan struct{}),
		pools:    pools,
	}

	go block.Listen()
	return block, nil
}

//Notify DiskStore about new Measurement
func (s *DiskStore) Notify(name string, m Measurement) {
	ownMeasurement := m.CopyFrom(s.pools)
	s.Add(name, m)
	s.pools.PutMeasurement(ownMeasurement)
}

//Add measurement to block
func (s *DiskStore) Add(name string, measurement Measurement) {
	doneChan := make(chan struct{})
	s.addChan <- addMessage{
		name:        name,
		doneChan:    doneChan,
		measurement: measurement,
	}
	<-doneChan
}

//Shutdown DiskBlock goroutine
func (s *DiskStore) Shutdown() {
	s.stopChan <- struct{}{}
}

//Listen for new measurements
func (s *DiskStore) Listen() {
	timer := time.NewTimer(10 * time.Second)
loop:
	for {
		timer.Stop()
		timer.Reset(time.Second * 5)
		select {
		case <-s.stopChan:
			s.commit()
			break loop
		case <-timer.C:
			s.commit()
		case message := <-s.addChan:
			s.handleAdd(message.name, message.measurement)
			message.doneChan <- struct{}{}
		}
	}
	s.cleanup()
}

func (s *DiskStore) cleanup() {
}

//Commit the buffered writes to actual disk
func (s *DiskStore) commit() {
	if s.block.Buffer.Len() == 0 {
		return
	}

	fileList, err := GetSortedFileList()
	if err != nil {
		fmt.Printf("couldn't get file List: %v", err)
	}
	defer s.block.Reset()
	if len(fileList) == 0 {
		WriteBlockToFile(s.block)
		return
	}
	latestFile := fileList[len(fileList)-1]
	if latestFile.size < maxFileSize {
		AppendBlockToFile(latestFile, s.block)
		return
	}
	WriteBlockToFile(s.block)

	if fileList.TotalSize() > maxDiskSize {
		oldestFile := fileList[0]
		os.Remove(filepath.Join(dataPath, oldestFile.name))
	}
}

func (s *DiskStore) handleAdd(name string, m Measurement) {
	id, err := s.meta.GetOrCreateID(name, m.Type())
	if err != nil {
		//measurement is probably of different type than it used to be, just ignore for now
		return
	}
	csvLineBytes, err := constructCsvLine(id, m)
	if err != nil {
		//ignore bad values
		return
	}
	s.block.AddBytes(m.Timestamp(), csvLineBytes)
	if s.block.Buffer.Len() > maxBuffer {
		s.commit()
	}

}
