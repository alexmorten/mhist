package mhist

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_timestampsFromFileName(t *testing.T) {
	Convey("gets timestamps from filename", t, func() {
		fileName := "1234-56789.csv"
		info, err := timestampsFromFileName(fileName)
		So(err, ShouldBeNil)
		So(info.oldestTs, ShouldEqual, 1234)
		So(info.latestTs, ShouldEqual, 56789)
	})
}
