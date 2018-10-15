package mhist_test

import (
	"testing"

	"github.com/codeuniversity/ppp-mhist"
	"github.com/codeuniversity/ppp-mhist/testhelpers"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStore(t *testing.T) {
	Convey("Store", t, func() {
		Convey("returns the correct map", func() {
			s := mhist.NewStore(100 * 1024 * 1024)
			for _, m := range testhelpers.GetSampleMeasurements(5, 1000, 20) {
				s.Add("temperature", m, false)
			}
			for _, m := range testhelpers.GetSampleMeasurements(6, 1040, 20) {
				s.Add("acceleration", m, false)
			}
			returnedMap := s.GetMeasurementsInTimeRange(1020, 1060, mhist.FilterDefinition{})
			So(len(returnedMap["temperature"]), ShouldEqual, 3)
			So(len(returnedMap["acceleration"]), ShouldEqual, 2)
			s.Shutdown()
		})
	})
}
