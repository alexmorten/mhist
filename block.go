package mhist

import (
	"math"
	"unsafe"
)

//Block keeps the timestamp range
type Block []SerializedMeasurement

//OldestTs ...
func (b Block) OldestTs() int64 {
	if len(b) == 0 {
		return 0
	}
	return b[0].Ts
}

//LatestTs ...
func (b Block) LatestTs() int64 {
	if len(b) == 0 {
		return 0
	}
	return b[len(b)-1].Ts

}

//Size of block
func (b Block) Size() int64 {
	return int64(len(b)) * serializedMeasurementSize
}

//UnderlyingByteSlice returns the memory representation of the Block
func (b Block) UnderlyingByteSlice() []byte {
	pointer := unsafe.Pointer(&b)

	pointerToRawBytes := (*(*[math.MaxInt32]byte))(pointer)
	byteSlice := (*(*pointerToRawBytes))[:b.Size()]
	return byteSlice
}

//BlockFromByteSlice interprets the byteSlice as a Block
//changing b after giving it to this func will result in unexpected changes to the returned Block
//you should discard the byteSlice after giving it to this function
func BlockFromByteSlice(b []byte) Block {
	amountOfMeasurements := len(b) / int(serializedMeasurementSize)
	pointer := unsafe.Pointer(&b)
	pointerToBlock := (*Block)(pointer)
	return (*pointerToBlock)[:amountOfMeasurements]
}
