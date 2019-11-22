package transform

import (
	"testing"

	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	companyName   = "companyName"
	companyNumber = "companyNumber"
	companyStatus = "companyStatus"
	companyType   = "companyType"

	id = "id"

	nameStart = "nameStart"
	nameEnd   = "nameEnd"
)

func TestUnitTransformMongoCompanyToEsCompany(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	mf := format.NewMockFormatter(ctrl)
	mwf := NewTransformer(mw, mf)

	Convey("Given I have a fully populated mongoCompany", t, func() {

		md := datastructures.MongoData{
			CompanyName:   companyName,
			CompanyNumber: companyNumber,
			CompanyStatus: companyStatus,
			CompanyType:   companyType,
			Links:         datastructures.MongoLinks{},
		}

		mc := datastructures.MongoCompany{
			ID:   id,
			Data: &md,
		}

		Convey("When I call TransformMongoCompanyToEsCompany", func() {

			mf.EXPECT().SplitCompanyNameEndings(md.CompanyName).Return(nameStart, nameEnd)

			esData := mwf.TransformMongoCompanyToEsCompany(&mc)

			Convey("Then I expect a fully populated EsItem", func() {

				So(esData, ShouldNotBeNil)
				So(esData.CompanyType, ShouldEqual, companyType)
				So(esData.ID, ShouldEqual, id)
				So(esData.Items.CompanyNumber, ShouldEqual, companyNumber)
				So(esData.Items.CompanyStatus, ShouldEqual, companyStatus)
				So(esData.Items.CorporateName, ShouldEqual, companyName)
				So(esData.Items.CorporateNameStart, ShouldEqual, nameStart)
				So(esData.Items.CorporateNameEnding, ShouldEqual, nameEnd)
			})
		})
	})

	Convey("Given my mongoCompany is not populated", t, func() {

		mc := datastructures.MongoCompany{}

		Convey("When I call TransformMongoCompanyToEsCompany", func() {

			esData := mwf.TransformMongoCompanyToEsCompany(&mc)

			Convey("I expect it to return nil", func() {
				So(esData, ShouldBeNil)
			})
		})
	})

	Convey("Given the companyName is empty", t, func() {

		md := datastructures.MongoData{
			CompanyName:   "",
			CompanyNumber: companyNumber,
			CompanyStatus: companyStatus,
			CompanyType:   companyType,
			Links:         datastructures.MongoLinks{},
		}

		mc := datastructures.MongoCompany{
			ID:   id,
			Data: &md,
		}

		Convey("Then I expect an error to be logged", func() {

			mw.EXPECT().LogMissingCompanyName(mc.ID).Times(1)

			Convey("When I call TransformMongoCompanyToEsCompany", func() {

				esData := mwf.TransformMongoCompanyToEsCompany(&mc)

				Convey("And I expect esData to be nil", func() {

					So(esData, ShouldBeNil)
				})
			})
		})
	})
}
