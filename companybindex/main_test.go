package main

import (
	"encoding/json"
	"fmt"
	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/eshttp"
	"github.com/companieshouse/elasticsearch-data-loader/transform"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestGetAlphaKeys(t *testing.T) {

	Convey("Should obtain keys", t, func() {

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		client := eshttp.NewMockClient(ctrl)
		companies := []*datastructures.MongoCompany{{}}

		companyNames := []datastructures.CompanyName{{Name: "Blah Co"}}
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

	Convey("Should handle failure to marshal company names by exiting program", t, func() {

		// Save current function and restore at the end.
		realMarshal := marshal
		defer func() { marshal = realMarshal }()
		marshal = func(v interface{}) ([]byte, error) {
			return nil, &json.UnsupportedValueError{
				Value: reflect.Value{},
				Str:   "Test generated error",
			}
		}

		realFatalf := fatalf
		defer func() { fatalf = realFatalf }()
		fatalf = func(format string, v ...interface{}) {
			errorMessage := fmt.Sprintf(format, v)
			// We replace os.Exit() with panic() because it too exits execution at the right point,
			// but the GoConvey test framework can detect the latter.
			panic(errorMessage)
		}

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		client := eshttp.NewMockClient(ctrl)
		companies := []*datastructures.MongoCompany{{}}

		companyNames := []datastructures.CompanyName{{Name: "@"}}
		companyNamesBody, _ := json.Marshal(companyNames)

		transformer.EXPECT().GetCompanyNames(&companies, 0).Times(1)
		client.EXPECT().GetAlphaKeys(companyNamesBody, alphakeyURL).Times(0)

		So(func() { getAlphaKeys(transformer, &companies, 0, client) },
			ShouldPanicWith,
			"error marshal to json: [json: unsupported value: Test generated error]")

	})

}
