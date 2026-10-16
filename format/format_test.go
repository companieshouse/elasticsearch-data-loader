package format

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitSplitCompanyNameEndings(t *testing.T) {
	f := NewFormatter()

	Convey("SplitCompanyNameEndings should split standard endings", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST LIMITED")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, " LIMITED")
	})

	Convey("SplitCompanyNameEndings should prefer the longest valid ending", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST PUBLIC LIMITED COMPANY")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, " PUBLIC LIMITED COMPANY")
	})

	Convey("SplitCompanyNameEndings should not treat partial token matches as endings", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST UNLTD")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, " UNLTD")
	})

	Convey("SplitCompanyNameEndings should preserve trailing punctuation on endings", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST PCC LTD...")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, " PCC LTD...")
	})

	Convey("SplitCompanyNameEndings should return original name when no ending is present", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST HOLDINGS")
		So(nameStart, ShouldEqual, "TEST HOLDINGS")
		So(nameEnd, ShouldEqual, "")
	})
}
