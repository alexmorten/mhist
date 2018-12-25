package mhist

import (
	"testing"

	"github.com/alexmorten/mhist/models"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_BlockMemoryDumpingAndLoading(t *testing.T) {
	Convey("Underlying bytes works", t, func() {
		m1 := SerializedMeasurement{ID: 1, Numerical: models.Numerical{Ts: 2, Value: 4}}
		m2 := SerializedMeasurement{ID: 5, Numerical: models.Numerical{Ts: 6, Value: 7}}
		m3 := SerializedMeasurement{ID: 100, Numerical: models.Numerical{Ts: 200, Value: 300}}
		block := Block{}
		block = append(block, m1)
		block = append(block, m2)
		block = append(block, m3)
		b := block.UnderlyingByteSlice()
		bCopy := make([]byte, len(b))
		copy(bCopy, b)

		newBlock := BlockFromByteSlice(bCopy)
		So(newBlock, ShouldResemble, Block{m1, m2, m3})
	})
}

var newBlock Block

func Benchmark_BlockMemoryDumpingAndLoading(b *testing.B) {
	block := Block{}
	for i := 0; i < b.N; i++ {
		m := SerializedMeasurement{ID: int64(i + 1), Numerical: models.Numerical{Ts: int64(i + 2), Value: float64(i) + 3.5}}
		block = append(block, m)
		byteSlice := block.UnderlyingByteSlice()
		newBlock = BlockFromByteSlice(byteSlice)
	}
}
