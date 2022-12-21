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

const (
	submitBulkToESCalled     = "When SubmitBulkToES is called"
	returnedBytesShouldBeNil = "And returnedBytes should be nil"
	errShouldNotBeNil        = "And err should not be nil"
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
	uri := esDestURL + "/" + esDestIndex + "/_bulk"

	Convey("Given a successful post of bulks to Elastic Search", t, func() {

		mr.EXPECT().Post(bulk, uri).Return(constructSuccessResponse(), nil)

		Convey(submitBulkToESCalled, func() {

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

		mr.EXPECT().Post(bulk, uri).Return(constructUnsuccessfulResponse(), errors.New("error posting bulk"))

		Convey("Then the post error should be logged", func() {

			mw.EXPECT().LogPostError(string(companyNumbers)).Times(1)

			Convey(submitBulkToESCalled, func() {

				returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

				Convey(returnedBytesShouldBeNil, func() {

					So(returnedBytes, ShouldBeNil)

					Convey(errShouldNotBeNil, func() {

						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})

	Convey("Given an unexpected response when posting bulks to Elastic Search", t, func() {

		mr.EXPECT().Post(bulk, uri).Return(constructUnsuccessfulResponse(), nil)

		Convey("Then the unexpected response should be logged", func() {

			mw.EXPECT().LogUnexpectedResponse(string(companyNumbers)).Times(1)

			Convey(submitBulkToESCalled, func() {

				returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

				Convey(returnedBytesShouldBeNil, func() {

					So(returnedBytes, ShouldBeNil)

					Convey(errShouldNotBeNil, func() {

						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}

func TestUnitGetAlphaKeys(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	mr := NewMockRequester(ctrl)
	mc := NewClientWithRequester(mw, mr)

	companyNames := make([]byte, 1)
	alphaKeyURL := "alphaKeyURL"
	uri := alphaKeyURL + "/alphakey-bulk"

	Convey("Given a successful post of company names to the alpha key service", t, func() {

		mr.EXPECT().Post(companyNames, uri).Return(constructSuccessResponse(), nil)

		Convey("When GetAlphaKeys is called", func() {

			returnedBytes, err := mc.GetAlphaKeys(companyNames, alphaKeyURL)

			Convey("Then returnedBytes should not be nil", func() {

				So(returnedBytes, ShouldNotBeNil)

				Convey("And err should be nil", func() {

					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given an unsuccessful post of company names to the alpha key service", t, func() {

		mr.EXPECT().Post(companyNames, uri).Return(constructUnsuccessfulResponse(), errors.New("error posting company names to the alpha key service"))

		Convey("Then the alpha ker error should be logged", func() {

			mw.EXPECT().LogAlphaKeyErrors(string(companyNames)).Times(1)

			Convey("When GetAlphaKeys is called", func() {

				returnedBytes, err := mc.GetAlphaKeys(companyNames, alphaKeyURL)

				Convey(returnedBytesShouldBeNil, func() {

					So(returnedBytes, ShouldBeNil)

					Convey(errShouldNotBeNil, func() {

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
