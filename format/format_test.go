package format

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitSplitCompanyNameEndings(t *testing.T) {

	f := NewFormatter()

	Convey("Given a company name of 'TEST LIMITED'", t, func() {

		coName := "TEST LIMITED"

		Convey("When SplitCompanyNameEndings is called", func() {

			nameStart, nameEnd := f.SplitCompanyNameEndings(coName)

			Convey("Then nameStart should equal 'TEST'", func() {

				So(nameStart, ShouldEqual, "TEST")

				Convey("And nameEnd should equal ' LIMITED'", func() {

					So(nameEnd, ShouldEqual, " LIMITED")
				})
			})
		})
	})
}
