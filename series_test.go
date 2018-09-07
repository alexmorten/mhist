package mhist_test

import (
	"testing"

	"github.com/codeuniversity/ppp-mhist"
	"github.com/codeuniversity/ppp-mhist/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSeries(t *testing.T) {
	Convey("Series", t, func() {
		Convey("returns no measurements if empty", func() {
			s := mhist.NewSeries()
			returnedMeasurements := s.GetMeasurementsInTimeRange(1005, 1035)
			s.Shutdown()

			So(len(returnedMeasurements), ShouldEqual, 0)
		})
		Convey("returns correct measurements if given range is inside", func() {
			s := mhist.NewSeries()
			test_helpers.AddMeasurementsToSeries(s)
			returnedMeasurements := s.GetMeasurementsInTimeRange(1005, 1035)

			s.Shutdown()
			So(len(returnedMeasurements), ShouldEqual, 3)
			So(returnedMeasurements[0].Value, ShouldEqual, 11)
			So(returnedMeasurements[1].Value, ShouldEqual, 12)
			So(returnedMeasurements[2].Value, ShouldEqual, 13)
		})
		Convey("returns all measurements if it is completly inside given range", func() {
			s := mhist.NewSeries()
			test_helpers.AddMeasurementsToSeries(s)
			returnedMeasurements := s.GetMeasurementsInTimeRange(500, 4000)

			s.Shutdown()
			So(len(returnedMeasurements), ShouldEqual, 5)
		})

		Convey("returns no measurements if given range has no overlap", func() {
			s := mhist.NewSeries()
			test_helpers.AddMeasurementsToSeries(s)
			returnedMeasurements := s.GetMeasurementsInTimeRange(3000, 4000)

			s.Shutdown()
			So(len(returnedMeasurements), ShouldEqual, 0)
		})

		Convey("returns correct if given range has partialy overlaps", func() {
			s := mhist.NewSeries()
			test_helpers.AddMeasurementsToSeries(s)
			returnedMeasurements := s.GetMeasurementsInTimeRange(1025, 4000)

			s.Shutdown()
			So(len(returnedMeasurements), ShouldEqual, 2)
			So(returnedMeasurements[0].Value, ShouldEqual, 13)
			So(returnedMeasurements[1].Value, ShouldEqual, 14)
		})
	})
}
