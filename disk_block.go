package mhist

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maxBuffer = 10 * 1024 * 1024

//DiskBlock handles buffered writes to Disk
type DiskBlock struct {
	writer          *bufio.Writer
	file            *os.File
	addChan         chan addMessage
	stopChan        chan struct{}
	measurementType MeasurementType
	needsSync       bool
	bufferedSize    int
}

type addMessage struct {
	measurement Measurement
	doneChan    chan struct{}
}

//NewDiskBlock initializes the DiskBlockRoutine
func NewDiskBlock(measurementType MeasurementType, name string) (*DiskBlock, error) {
	file, err := os.OpenFile(filePathTo(name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	block := &DiskBlock{
		writer:          bufio.NewWriter(file),
		file:            file,
		addChan:         make(chan addMessage),
		stopChan:        make(chan struct{}),
		measurementType: measurementType,
	}
	go block.Listen()
	return block, nil
}

//Add measurement to block
func (b *DiskBlock) Add(measurement Measurement) {
	doneChan := make(chan struct{})
	b.addChan <- addMessage{
		doneChan:    doneChan,
		measurement: measurement,
	}
	<-doneChan
}

//Shutdown DiskBlock goroutine
func (b *DiskBlock) Shutdown() {
	b.stopChan <- struct{}{}
}

//Listen for new measurements
func (b *DiskBlock) Listen() {
	timer := time.NewTimer(10 * time.Second)
loop:
	for {
		timer.Stop()
		timer.Reset(time.Second * 5)
		select {
		case <-b.stopChan:
			b.commit()
			break loop
		case <-timer.C:
			b.commit()
		case message := <-b.addChan:
			b.handleAdd(message.measurement)
			message.doneChan <- struct{}{}
		}
	}
	b.cleanup()
}

func (b *DiskBlock) cleanup() {
	b.writer.Flush()
	b.file.Close()
}

//Commit the buffered writes to actual disk
func (b *DiskBlock) commit() {
	if b.needsSync {
		b.writer.Flush()
		b.needsSync = false
		b.bufferedSize = 0
	}
}

func (b *DiskBlock) handleAdd(m Measurement) {
	if b.measurementType == m.Type() {
		byteSlice, err := constructCsvLine(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		b.writer.Write(byteSlice)
		b.bufferedSize += len(byteSlice)
		b.needsSync = true

		if b.bufferedSize > maxBuffer {
			b.commit()
		}

		return
	}
	fmt.Println(m, " is not the correct type for this series")
}

func filePathTo(name string) string {
	return filepath.Join(dataPath, name)
}
