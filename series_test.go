package mhist_test

import (
	"testing"

	"github.com/codeuniversity/ppp-mhist"
	"github.com/codeuniversity/ppp-mhist/testhelpers"
	. "github.com/smartystreets/goconvey/convey"
)

const maxSize = 100 * 1024 * 1024

func TestSeries(t *testing.T) {
	Convey("Series", t, func() {
		Convey("GetMeasurementsInTimeRange()", func() {
			Convey("returns no measurements if empty", func() {
				s := mhist.NewSeries(maxSize)
				returnedMeasurements := s.GetMeasurementsInTimeRange(1005, 1035)
				s.Shutdown()

				So(len(returnedMeasurements), ShouldEqual, 0)
			})
			Convey("returns correct measurements if given range is inside", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements := s.GetMeasurementsInTimeRange(1005, 1035)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 3)
			})
			Convey("returns all measurements if it is completly inside given range", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements := s.GetMeasurementsInTimeRange(500, 4000)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 5)
			})

			Convey("returns no measurements if given range has no overlap", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements := s.GetMeasurementsInTimeRange(3000, 4000)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 0)
			})

			Convey("returns correct if given range has partialy overlaps", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements := s.GetMeasurementsInTimeRange(1025, 4000)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 2)
			})
		})

		Convey("CutoffBelow()", func() {
			Convey("returns correct measurements", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)

				So(s.Size(), ShouldEqual, 80)
				returnedMeasurements := s.CutoffBelow(1025)
				So(len(returnedMeasurements), ShouldEqual, 3)
				So(s.Size(), ShouldEqual, 32)
				s.Shutdown()
			})

			Convey("returns no measurements if timestamp is below all of series", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)

				So(s.Size(), ShouldEqual, 80)
				returnedMeasurements := s.CutoffBelow(900)
				So(len(returnedMeasurements), ShouldEqual, 0)
				So(s.Size(), ShouldEqual, 80)
				s.Shutdown()
			})

			Convey("returns all measurements if timestamp is above all of series", func() {
				s := mhist.NewSeries(maxSize)
				testhelpers.AddMeasurementsToSeries(s)

				So(s.Size(), ShouldEqual, 80)
				returnedMeasurements := s.CutoffBelow(2000)
				So(len(returnedMeasurements), ShouldEqual, 5)
				So(s.Size(), ShouldEqual, 0)
				s.Shutdown()
			})
		})
	})
}
