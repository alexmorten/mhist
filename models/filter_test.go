package models

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Passes(t *testing.T) {
	Convey("correct timestamps pass the filter", t, func() {
		definition := FilterDefinition{
			Granularity: 2 * time.Millisecond,
			Names:       []string{"bla", "blup"},
		}
		filter := NewFilterCollection(definition)
		So(filter.Passes("foo", &Numerical{Ts: 1000000}), ShouldBeFalse)
		So(filter.Passes("bla", &Numerical{Ts: 1000000}), ShouldBeTrue)
		So(filter.Passes("bla", &Numerical{Ts: 2000000}), ShouldBeFalse)
		So(filter.Passes("bla", &Numerical{Ts: 3000000}), ShouldBeTrue)
		So(filter.Passes("bla", &Numerical{Ts: 4000000}), ShouldBeFalse)
	})
}

func Test_TimestampFilter_Passes(t *testing.T) {
	Convey("correct timestamps pass the filter", t, func() {
		filter := &TimestampFilter{Granularity: 2 * time.Millisecond}
		So(filter.Passes(&Numerical{Ts: 1000000}), ShouldBeTrue)
		So(filter.Passes(&Numerical{Ts: 2000000}), ShouldBeFalse)
		So(filter.Passes(&Numerical{Ts: 3000000}), ShouldBeTrue)
		So(filter.Passes(&Numerical{Ts: 4000000}), ShouldBeFalse)
	})
}
