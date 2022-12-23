package main

import (
	"encoding/json"
	"errors"
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

		// Save real functions so that we restore them at the end.
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
			// but the GoConvey test framework can detect the latter only.
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

	Convey("Should handle failure to get alpha keys by exiting program", t, func() {

		// Save real functions so that we restore them at the end.
		realFatalf := fatalf
		defer func() { fatalf = realFatalf }()
		fatalf = func(format string, v ...interface{}) {
			errorMessage := fmt.Sprintf(format, v)
			// We replace os.Exit() with panic() because it too exits execution at the right point,
			// but the GoConvey test framework can detect the latter only.
			panic(errorMessage)
		}

		unmarshalCalled := false
		realUnmarshal := unmarshal
		defer func() { unmarshal = realUnmarshal }()
		unmarshal = func(data []byte, v interface{}) error {
			unmarshalCalled = true
			return nil
		}

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		client := eshttp.NewMockClient(ctrl)
		companies := []*datastructures.MongoCompany{{}}

		companyNames := []datastructures.CompanyName{{Name: "@"}}
		companyNamesBody, _ := json.Marshal(companyNames)

		transformer.EXPECT().GetCompanyNames(&companies, 0).Return(companyNames)
		client.EXPECT().GetAlphaKeys(companyNamesBody, alphakeyURL).Return(
			nil, errors.New("Test generated error"))

		So(func() { getAlphaKeys(transformer, &companies, 0, client) },
			ShouldPanicWith,
			"error fetching alpha keys: [Test generated error]")
		So(unmarshalCalled, ShouldBeFalse)

	})

	Convey("Should handle failure to unmarshal company names by exiting program", t, func() {

		// Save real functions so that we restore them at the end.
		realUnmarshal := unmarshal
		defer func() { unmarshal = realUnmarshal }()
		unmarshal = func(data []byte, v interface{}) error {
			return &json.UnmarshalTypeError{
				Value:  "Test generated error",
				Type:   reflect.TypeOf(""),
				Offset: 0,
				Struct: "struct",
				Field:  "field",
			}
		}

		realFatalf := fatalf
		defer func() { fatalf = realFatalf }()
		fatalf = func(format string, v ...interface{}) {
			errorMessage := fmt.Sprintf(format, v)
			// We replace os.Exit() with panic() because it too exits execution at the right point,
			// but the GoConvey test framework can detect the latter only.
			panic(errorMessage)
		}

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		client := eshttp.NewMockClient(ctrl)
		companies := []*datastructures.MongoCompany{{}}

		companyNames := []datastructures.CompanyName{{Name: "@"}}
		companyNamesBody, _ := json.Marshal(companyNames)

		transformer.EXPECT().GetCompanyNames(&companies, 0).Return(companyNames)
		client.EXPECT().GetAlphaKeys(companyNamesBody, alphakeyURL).Return(
			[]byte("[{\"sameAsAlphaKey\":\"true\", \"orderedAlphaKey\":\"blah\"}]"), nil)

		So(func() { getAlphaKeys(transformer, &companies, 0, client) },
			ShouldPanicWith,
			"error [json: cannot unmarshal Test generated error into Go struct field "+
				"struct.field of type string [91 123 34 110 97 109 101 34 58 34 64 34 125 93]] "+
				"unmarshalling alphakey response for %!s(MISSING)")

	})

}
