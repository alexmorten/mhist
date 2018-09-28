package mhist

import (
	"bytes"
)

//Block keeps the timestamp range
type Block struct {
	Buffer          bytes.Buffer
	oldestTimestamp int64
	latestTimestamp int64
}

//AddBytes to block
func (b *Block) AddBytes(ts int64, byteSlice []byte) {
	_, err := b.Buffer.Write(byteSlice)
	if err != nil {
		return
	}

	if b.oldestTimestamp == 0 {
		b.oldestTimestamp = ts
	}

	b.latestTimestamp = ts

}

//OldestTs ...
func (b *Block) OldestTs() int64 {
	return b.oldestTimestamp
}

//LatestTs ...
func (b *Block) LatestTs() int64 {
	return b.latestTimestamp
}

//Reset Block, i.e. after writing
func (b *Block) Reset() {
	b.Buffer.Reset()
	b.oldestTimestamp = 0
	b.latestTimestamp = 0
}
