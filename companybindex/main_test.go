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

func TestUnitGetAlphaKeys(t *testing.T) {

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

		restoreJsonMarshal := stubJsonMarshal()
		defer restoreJsonMarshal()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

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
			"error marshal to json: json: unsupported value: Test generated error")

	})

	Convey("Should handle failure to get alpha keys by exiting program", t, func() {

		unmarshalCalled := false
		restoreJsonUnmarshal := mockJsonUnmarshal(&unmarshalCalled)
		defer restoreJsonUnmarshal()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

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
			"error fetching alpha keys: Test generated error")
		So(unmarshalCalled, ShouldBeFalse)

	})

	Convey("Should handle failure to unmarshal company names by exiting program", t, func() {

		restoreJsonUnmarshal := stubJsonUnmarshalWithError()
		defer restoreJsonUnmarshal()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

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
			"error json: cannot unmarshal Test generated error into Go struct "+
				"field struct.field of type string unmarshalling alphakey response for"+
				" [{\"name\":\"@\"}]")

	})

}

func TestUnitTransformMongoCompaniesToEsCompanies(t *testing.T) {

	Convey("Should transform mongo companies to elasticsearch companies", t, func() {

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		companies := []*datastructures.MongoCompany{{
			ID: "Co",
		}}
		keys := []datastructures.AlphaKey{{
			SameAsAlphaKey:  "true",
			OrderedAlphaKey: "blah",
		}}

		transformer.EXPECT().TransformMongoCompanyToEsCompany(
			&datastructures.MongoCompany{
				ID: "Co",
			},
			&datastructures.AlphaKey{
				SameAsAlphaKey:  "true",
				OrderedAlphaKey: "blah",
			}).Return(&datastructures.EsCompany{
			ID:                    "",
			CompanyType:           "",
			Items:                 datastructures.EsItem{},
			Kind:                  "",
			Links:                 nil,
			OrderedAlphaKeyWithID: "",
		})

		bulk, companyNumbers, target :=
			transformMongoCompaniesToEsCompanies(
				1,
				transformer,
				&companies,
				keys,
				[]byte("bulk"),
				[]byte("companyNumbers"),
				1)
		So(string(bulk), ShouldContainSubstring,
			`bulk{ "create": { "_id": "" } }
{"ID":"","company_type":"","items":{"company_number":"","corporate_name":"","corporate_name_start":`+
				`"","record_type":"","alpha_key":"","ordered_alpha_key":""},`+
				`"kind":"","links":null,"ordered_alpha_key_with_id":""}`)
		So(string(companyNumbers), ShouldEqual, `companyNumbers
`)
		So(target, ShouldEqual, 1)
	})

	Convey("Should handle failure to marshal company by exiting program", t, func() {

		restoreJsonMarshal := stubJsonMarshal()
		defer restoreJsonMarshal()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		companies := []*datastructures.MongoCompany{{
			ID: "Co",
		}}
		keys := []datastructures.AlphaKey{{
			SameAsAlphaKey:  "true",
			OrderedAlphaKey: "blah",
		}}

		transformer.EXPECT().TransformMongoCompanyToEsCompany(
			&datastructures.MongoCompany{
				ID: "Co",
			},
			&datastructures.AlphaKey{
				SameAsAlphaKey:  "true",
				OrderedAlphaKey: "blah",
			}).Return(&datastructures.EsCompany{
			ID:                    "",
			CompanyType:           "",
			Items:                 datastructures.EsItem{},
			Kind:                  "",
			Links:                 nil,
			OrderedAlphaKeyWithID: "",
		})

		So(func() {
			transformMongoCompaniesToEsCompanies(
				1,
				transformer,
				&companies,
				keys,
				[]byte("bulk"),
				[]byte("companyNumbers"),
				1)
		},
			ShouldPanicWith,
			"error marshal to json: json: unsupported value: Test generated error")

	})

	Convey("Should increment skip count where company is nil", t, func() {

		restoreSkipChannel := stubSkipChannel()
		defer restoreSkipChannel()

		ctrl := gomock.NewController(t)
		transformer := transform.NewMockTransformer(ctrl)
		companies := []*datastructures.MongoCompany{{
			ID: "Co",
		}}
		keys := []datastructures.AlphaKey{{
			SameAsAlphaKey:  "true",
			OrderedAlphaKey: "blah",
		}}

		transformer.EXPECT().TransformMongoCompanyToEsCompany(
			&datastructures.MongoCompany{
				ID: "Co",
			},
			&datastructures.AlphaKey{
				SameAsAlphaKey:  "true",
				OrderedAlphaKey: "blah",
			}).Return(nil)

		go transformMongoCompaniesToEsCompanies(
			1,
			transformer,
			&companies,
			keys,
			[]byte("bulk"),
			[]byte("companyNumbers"),
			1)
		increment := <-skipChannel
		So(increment, ShouldEqual, 1)
	})
}

func TestUnitSubmitBulkToES(t *testing.T) {

	Convey("Should report bulk submission success", t, func() {

		unmarshalCalled := false
		restoreJsonUnmarshal := mockJsonUnmarshal(&unmarshalCalled)
		defer restoreJsonUnmarshal()

		ctrl := gomock.NewController(t)
		client := eshttp.NewMockClient(ctrl)

		client.EXPECT().SubmitBulkToES([]byte("bulk"), []byte("companyNumbers"), esDestURL, esDestIndex).
			Return([]byte("bulk"), nil)

		submissionFailed := submitBulkToES(nil, client, []byte("bulk"), []byte("companyNumbers"))
		So(submissionFailed, ShouldEqual, false)
	})

	Convey("Should report bulk submission failure", t, func() {

		ctrl := gomock.NewController(t)
		client := eshttp.NewMockClient(ctrl)

		client.EXPECT().SubmitBulkToES([]byte("bulk"), []byte("companyNumbers"), esDestURL, esDestIndex).
			Return([]byte("bulk"), errors.New("Test generated error"))

		submissionFailed := submitBulkToES(nil, client, []byte("bulk"), []byte("companyNumbers"))
		So(submissionFailed, ShouldEqual, true)
	})

	Convey("Should handle failure to unmarshal bulk response by exiting program", t, func() {

		restoreJsonUnmarshal := stubJsonUnmarshalWithError()
		defer restoreJsonUnmarshal()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		ctrl := gomock.NewController(t)
		client := eshttp.NewMockClient(ctrl)

		client.EXPECT().SubmitBulkToES([]byte("bulk"), []byte("companyNumbers"), esDestURL, esDestIndex).
			Return([]byte("bulk"), nil)

		So(func() {
			submitBulkToES(nil, client, []byte("bulk"), []byte("companyNumbers"))
		},
			ShouldPanicWith,
			"error unmarshalling json: [json: cannot unmarshal Test generated error into Go "+
				"struct field struct.field of type string] actual response: [bulk]")

	})

	Convey("Should handle failure to create elasticsearch document by exiting program", t, func() {

		restoreJsonUnmarshal := stubJsonUnmarshalWithEsDocumentCreationResponseError()
		defer restoreJsonUnmarshal()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		ctrl := gomock.NewController(t)
		client := eshttp.NewMockClient(ctrl)

		client.EXPECT().SubmitBulkToES([]byte("bulk"), []byte("companyNumbers"), esDestURL, esDestIndex).
			Return([]byte("bulk"), nil)

		So(func() {
			submitBulkToES(nil, client, []byte("bulk"), []byte("companyNumbers"))
		},
			ShouldPanicWith,
			"error inserting doc: Test generated error")

	})
}

func stubJsonMarshal() func() {
	// Stub out json.Marshal
	realMarshal := marshal
	marshal = func(v interface{}) ([]byte, error) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.Value{},
			Str:   "Test generated error",
		}
	}
	// Return function to restore json.Marshal
	return func() { marshal = realMarshal }
}

func stubLogFatalf() func() {
	// Stub out log.Fatalf
	realFatalf := fatalf
	fatalf = func(format string, v ...interface{}) {
		errorMessage := fmt.Sprintf(format, v...)
		// We replace os.Exit() with panic() because it too exits execution at the right point,
		// but the GoConvey test framework can detect the latter only.
		panic(errorMessage)
	}
	// Return function to restore log.Fatalf
	return func() { fatalf = realFatalf }
}

func mockJsonUnmarshal(unmarshalCalled *bool) func() {
	// Mock out json.Unmarshal
	realUnmarshal := unmarshal
	unmarshal = func(data []byte, v interface{}) error {
		*unmarshalCalled = true
		return nil
	}
	// Return function to restore json.Unmarshal
	return func() { unmarshal = realUnmarshal }
}

func stubJsonUnmarshalWithError() func() {
	// Stub out json.Unmarshal
	realUnmarshal := unmarshal
	unmarshal = func(data []byte, v interface{}) error {
		return &json.UnmarshalTypeError{
			Value:  "Test generated error",
			Type:   reflect.TypeOf(""),
			Offset: 0,
			Struct: "struct",
			Field:  "field",
		}
	}
	// Return function to restore json.Unmarshal
	return func() { unmarshal = realUnmarshal }
}

func stubJsonUnmarshalWithEsDocumentCreationResponseError() func() {
	// Stub out json.Unmarshal
	realUnmarshal := unmarshal
	unmarshal = func(data []byte, v interface{}) error {
		bulkResponse := v.(*esBulkResponse)
		bulkResponse.Errors = true
		bulkResponse.Items = make([]esBulkItemResponse, 1)
		bulkResponse.Items[0] =
			map[string]esBulkItemResponseData{
				"create": {Index: "Index", ID: "Id", Status: 500, Error: "Test generated error"},
			}
		return nil
	}
	// Return function to restore json.Unmarshal
	return func() { unmarshal = realUnmarshal }
}

func stubSkipChannel() func() {
	// Stub out skipChannel
	realSkipChannel := skipChannel
	skipChannel = make(chan int)

	// Return function to restore skipChannel
	return func() { skipChannel = realSkipChannel }
}
