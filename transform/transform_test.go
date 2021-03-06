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

	sameAsAlphaKey  = "sameAsAlphaKey"
	orderedAlphaKey = "orderedAlphaKey"

	companyOne   = "companyOne"
	companyTwo   = "companyTwo"
	companyThree = "companyThree"
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

		ak := datastructures.AlphaKey{
			SameAsAlphaKey:  sameAsAlphaKey,
			OrderedAlphaKey: orderedAlphaKey,
		}

		Convey("When I call TransformMongoCompanyToEsCompany", func() {

			mf.EXPECT().SplitCompanyNameEndings(md.CompanyName).Return(nameStart, nameEnd)

			esData := mwf.TransformMongoCompanyToEsCompany(&mc, &ak)

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

		ak := datastructures.AlphaKey{}

		Convey("When I call TransformMongoCompanyToEsCompany", func() {

			esData := mwf.TransformMongoCompanyToEsCompany(&mc, &ak)

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

		ak := datastructures.AlphaKey{}

		Convey("Then I expect an error to be logged", func() {

			mw.EXPECT().LogMissingCompanyName(mc.ID).Times(1)

			Convey("When I call TransformMongoCompanyToEsCompany", func() {

				esData := mwf.TransformMongoCompanyToEsCompany(&mc, &ak)

				Convey("And I expect esData to be nil", func() {

					So(esData, ShouldBeNil)
				})
			})
		})
	})
}

func TestUnitGetCompanyNames(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	mf := format.NewMockFormatter(ctrl)
	mwf := NewTransformer(mw, mf)

	Convey("Given I have an array of three mongo companies", t, func() {

		mc1 := datastructures.MongoCompany{
			Data: &datastructures.MongoData{
				CompanyName: companyOne,
			},
		}

		mc2 := datastructures.MongoCompany{
			Data: &datastructures.MongoData{
				CompanyName: companyTwo,
			},
		}

		mc3 := datastructures.MongoCompany{
			Data: &datastructures.MongoData{
				CompanyName: companyThree,
			},
		}

		companies := []*datastructures.MongoCompany{&mc1, &mc2, &mc3}

		Convey("When I call GetCompanyNames", func() {

			companyNames := mwf.GetCompanyNames(&companies, 3)

			Convey("Then I expect 3 CompanyNames to be returned", func() {

				So(len(companyNames), ShouldEqual, 3)

				Convey("And the names should be in order", func() {

					So(companyNames[0].Name, ShouldEqual, companyOne)
					So(companyNames[1].Name, ShouldEqual, companyTwo)
					So(companyNames[2].Name, ShouldEqual, companyThree)
				})
			})
		})
	})
}
