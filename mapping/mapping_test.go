package mapping

import (
	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMapResult(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWrite(ctrl)
	mf := format.NewMockFormatter(ctrl)
	mwf := Mapper{Writer: mw, Format: mf}

	Convey("Given I have a fully populated mongoCompany", t, func() {

		md := datastructures.MongoData{
			CompanyName:   "EXAMPLE LIMITED",
			CompanyNumber: "45454",
			CompanyStatus: "active",
			CompanyType:   "limited",
			Links:         datastructures.MongoLinks{},
		}

		mc := datastructures.MongoCompany{
			ID:   "565656",
			Data: &md,
		}

		Convey("When I call mapResult", func() {

			mf.EXPECT().SplitCompanyNameEndings(md.CompanyName).Return("foo", "bar")

			esData := mwf.MapResult(&mc)

			Convey("Then I expect a fully populated EsItem", func() {
				So(esData, ShouldNotBeNil)
			})
		})
	})

	Convey("Given my mongoCompany is not populated", t, func() {

		md := datastructures.MongoData{}
		mc := datastructures.MongoCompany{}

		Convey("When I call mapResult", func() {

			mf.EXPECT().SplitCompanyNameEndings(md.CompanyName).Return("foo", "bar")

			esData := mwf.MapResult(&mc)

			Convey("I expect it to return nil", func() {
				So(esData, ShouldBeNil)
			})
		})
	})

	Convey("Given the companyName is empty", t, func() {
		md := datastructures.MongoData{
			CompanyName:   "",
			CompanyNumber: "45454",
			CompanyStatus: "active",
			CompanyType:   "limited",
			Links:         datastructures.MongoLinks{},
		}

		mc := datastructures.MongoCompany{
			ID:   "565656",
			Data: &md,
		}

		Convey("When I call mapResult", func() {

			mw.EXPECT().WriteToFile3(mc.ID).String()
			mf.EXPECT().SplitCompanyNameEndings(md.CompanyName).Return("foo", "bar")
			esData := mwf.MapResult(&mc)

			Convey("I expect it to return nil", func() {

				So(esData, ShouldBeNil)
			})
		})
	})
}
