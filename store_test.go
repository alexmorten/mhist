package mhist_test

import (
	"testing"

	"github.com/codeuniversity/ppp-mhist"
	"github.com/codeuniversity/ppp-mhist/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStore(t *testing.T) {
	Convey("Store", t, func() {
		Convey("returns the correct map", func() {
			s := mhist.NewStore()
			for _, m := range test_helpers.GetSampleMeasurements(5, 1000, 20) {
				s.Add("temperature", m)
			}
			for _, m := range test_helpers.GetSampleMeasurements(6, 1040, 20) {
				s.Add("acceleration", m)
			}
			returnedMap := s.GetAllMeasurementsInTimeRange(1020, 1060)
			So(len(returnedMap["temperature"]), ShouldEqual, 3)
			So(len(returnedMap["acceleration"]), ShouldEqual, 2)
			s.Shutdown()
		})
	})
}
