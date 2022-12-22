package main

import (
	"encoding/json"
	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/eshttp"
	"github.com/companieshouse/elasticsearch-data-loader/transform"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetAlphaKeys(t *testing.T) {

	Convey("Should obtain keys", t, func() {

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		client := eshttp.NewMockClient(ctrl)
		companies := []*datastructures.MongoCompany{{}}

		companyNames := []datastructures.CompanyName{{"Blah Co"}}
		companyNamesBody, _ := json.Marshal(companyNames)

		transformer.EXPECT().GetCompanyNames(&companies, 0).Return(companyNames)
		client.EXPECT().GetAlphaKeys(companyNamesBody, alphakeyURL).Return(
			[]byte("[{\"sameAsAlphaKey\":\"true\", \"orderedAlphaKey\":\"blah\"}]"), nil)
		expectedAlphaKey := datastructures.AlphaKey{
			SameAsAlphaKey:  "true",
			OrderedAlphaKey: "blah",
		}

		err, alphaKeys := getAlphaKeys(transformer, &companies, 0, client)
		So(err, ShouldBeNil)
		So(alphaKeys, ShouldNotBeNil)
		So(len(alphaKeys), ShouldEqual, 1)
		So(alphaKeys[0], ShouldResemble, expectedAlphaKey)
	})

}
