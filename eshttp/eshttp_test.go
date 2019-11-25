package eshttp

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/companieshouse/elasticsearch-data-loader/write"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitSubmitBulkToES(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	mr := NewMockRequester(ctrl)
	mc := NewClientWithRequester(mw, mr)

	bulk := make([]byte, 1)
	companyNumbers := make([]byte, 1)
	esDestURL := "esDestURL"
	esDestIndex := "esDestIndex"

	Convey("Given a successful post of bulks to Elastic Search", t, func() {

		mr.EXPECT().PostBulkToElasticSearch(bulk, esDestURL, esDestIndex).Return(constructSuccessResponse(), nil)

		Convey("When SubmitBulkToES is called", func() {

			returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

			Convey("Then returnedBytes should not be nil", func() {

				So(returnedBytes, ShouldNotBeNil)

				Convey("And err should be nil", func() {

					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given an unsuccessful post of bulks to Elastic Search", t, func() {

		mr.EXPECT().PostBulkToElasticSearch(bulk, esDestURL, esDestIndex).Return(constructUnsuccessfulResponse(), errors.New("error posting bulk"))

		Convey("Then the post error should be logged", func() {

			mw.EXPECT().LogPostError(string(companyNumbers)).Times(1)

			Convey("When SubmitBulkToES is called", func() {

				returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

				Convey("And returnedBytes should be nil", func() {

					So(returnedBytes, ShouldBeNil)

					Convey("And err should not be nil", func() {

						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})

	Convey("Given an unexpected response when posting bulks to Elastic Search", t, func() {

		mr.EXPECT().PostBulkToElasticSearch(bulk, esDestURL, esDestIndex).Return(constructUnsuccessfulResponse(), nil)

		Convey("Then the unexpected response should be logged", func() {

			mw.EXPECT().LogUnexpectedResponse(string(companyNumbers)).Times(1)

			Convey("When SubmitBulkToES is called", func() {

				returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

				Convey("And returnedBytes should be nil", func() {

					So(returnedBytes, ShouldBeNil)

					Convey("And err should not be nil", func() {

						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}

func constructSuccessResponse() *http.Response {

	return &http.Response{
		StatusCode: 201,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`Created`)),
		Header:     make(http.Header),
	}
}

func constructUnsuccessfulResponse() *http.Response {

	return &http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`Internal server error`)),
		Header:     make(http.Header),
	}
}
