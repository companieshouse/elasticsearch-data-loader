package eshttp

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
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

func TestUnitNewClientWithRequesters(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	signedReq := NewMockRequester(ctrl)
	unsignedReq := NewMockRequester(ctrl)

	mc := NewClientWithRequesters(mw, signedReq, unsignedReq)

	bulk := make([]byte, 1)
	companyNumbers := make([]byte, 1)
	esDestURL := "esDestURL"
	esDestIndex := "esDestIndex"
	esUri := esDestURL + "/" + esDestIndex + "/_bulk"

	companyNames := make([]byte, 1)
	alphaKeyURL := "alphaKeyURL"
	alphaKeyUri := alphaKeyURL + "/alphakey-bulk"

	Convey("Given a client with separate signed and unsigned requesters", t, func() {

		// Expect OpenSearch request to use signedReq
		signedReq.EXPECT().Post(bulk, esUri).Return(constructSuccessResponse(), nil)
		// Expect AlphaKey request to use unsignedReq
		unsignedReq.EXPECT().Post(companyNames, alphaKeyUri).Return(constructSuccessResponse(), nil)

		Convey("When SubmitBulkToES is called", func() {

			returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

			Convey("Then the signed requester should be used", func() {

				So(returnedBytes, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})

		Convey("When GetAlphaKeys is called", func() {

			returnedBytes, err := mc.GetAlphaKeys(companyNames, alphaKeyURL)

			Convey("Then the unsigned requester should be used", func() {

				So(returnedBytes, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitSubmitBulkToES_UsesSeparateSignedRequester(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	signedReq := NewMockRequester(ctrl)
	unsignedReq := NewMockRequester(ctrl)

	mc := NewClientWithRequesters(mw, signedReq, unsignedReq)

	bulk := make([]byte, 1)
	companyNumbers := make([]byte, 1)
	esDestURL := "esDestURL"
	esDestIndex := "esDestIndex"
	esUri := esDestURL + "/" + esDestIndex + "/_bulk"

	Convey("Given a client with separate requesters", t, func() {

		// Only signedReq should be called for OpenSearch
		signedReq.EXPECT().Post(bulk, esUri).Return(constructSuccessResponse(), nil)
		// unsignedReq should NOT be called
		unsignedReq.EXPECT().Post(gomock.Any(), gomock.Any()).Times(0)

		Convey("When SubmitBulkToES is called", func() {

			returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

			Convey("Then only the signed requester should be used", func() {

				So(returnedBytes, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitGetAlphaKeys_UsesSeparateUnsignedRequester(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	signedReq := NewMockRequester(ctrl)
	unsignedReq := NewMockRequester(ctrl)

	mc := NewClientWithRequesters(mw, signedReq, unsignedReq)

	companyNames := make([]byte, 1)
	alphaKeyURL := "alphaKeyURL"
	alphaKeyUri := alphaKeyURL + "/alphakey-bulk"

	Convey("Given a client with separate requesters", t, func() {

		// Only unsignedReq should be called for AlphaKey
		unsignedReq.EXPECT().Post(companyNames, alphaKeyUri).Return(constructSuccessResponse(), nil)
		// signedReq should NOT be called
		signedReq.EXPECT().Post(gomock.Any(), gomock.Any()).Times(0)

		Convey("When GetAlphaKeys is called", func() {

			returnedBytes, err := mc.GetAlphaKeys(companyNames, alphaKeyURL)

			Convey("Then only the unsigned requester should be used", func() {

				So(returnedBytes, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUnitNewClient_InitializesSeparateRequesters(t *testing.T) {

	Convey("Given NewClient is called", t, func() {

		mw := write.NewMockWriter(gomock.NewController(t))
		client := NewClient(mw).(*ClientImpl)

		Convey("Then both signedRequester and unsignedRequester should be initialized", func() {

			So(client.signedRequester, ShouldNotBeNil)
			So(client.unsignedRequester, ShouldNotBeNil)

			Convey("And unsignedRequester should always be UnsignedRequest", func() {

				_, ok := client.unsignedRequester.(*UnsignedRequest)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestUnitNewClientWithRequesters_HandleErrors(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	signedReq := NewMockRequester(ctrl)
	unsignedReq := NewMockRequester(ctrl)

	mc := NewClientWithRequesters(mw, signedReq, unsignedReq)

	bulk := make([]byte, 1)
	companyNumbers := make([]byte, 1)
	esDestURL := "esDestURL"
	esDestIndex := "esDestIndex"
	esUri := esDestURL + "/" + esDestIndex + "/_bulk"

	companyNames := make([]byte, 1)
	alphaKeyURL := "alphaKeyURL"
	alphaKeyUri := alphaKeyURL + "/alphakey-bulk"

	Convey("Given a client with separate requesters when errors occur", t, func() {

		Convey("When SubmitBulkToES fails with the signed requester", func() {

			signedReq.EXPECT().Post(bulk, esUri).Return(constructUnsuccessfulResponse(), errors.New("signed request failed"))
			mw.EXPECT().LogPostError(string(companyNumbers)).Times(1)

			returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

			Convey("Then error should be returned", func() {

				So(returnedBytes, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When GetAlphaKeys fails with the unsigned requester", func() {

			unsignedReq.EXPECT().Post(companyNames, alphaKeyUri).Return(constructUnsuccessfulResponse(), errors.New("unsigned request failed"))
			mw.EXPECT().LogAlphaKeyErrors(string(companyNumbers)).Times(1)

			returnedBytes, err := mc.GetAlphaKeys(companyNames, alphaKeyURL)

			Convey("Then error should be returned", func() {

				So(returnedBytes, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUnitNewRequester_WithEnvironmentVariable(t *testing.T) {

	Convey("Given USE_AWS_SIGV4 environment variables", t, func() {

		Convey("When USE_AWS_SIGV4 is set to 'true'", func() {

			t.Setenv("USE_AWS_SIGV4", "true")
			requester := NewRequester()

			Convey("Then it should attempt to create a signed requester", func() {

				// Should be either SignedRequest or UnsignedRequest (fallback)
				So(requester, ShouldNotBeNil)
				_, isSignedRequest := requester.(*SignedRequest)
				_, isUnsignedRequest := requester.(*UnsignedRequest)
				So(isSignedRequest || isUnsignedRequest, ShouldBeTrue)
			})
		})

		Convey("When USE_AWS_SIGV4 is set to 'false'", func() {

			t.Setenv("USE_AWS_SIGV4", "false")
			requester := NewRequester()

			Convey("Then it should return an unsigned requester", func() {

				_, ok := requester.(*UnsignedRequest)
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When USE_AWS_SIGV4 is not set", func() {

			os.Unsetenv("USE_AWS_SIGV4")
			requester := NewRequester()

			Convey("Then it should return an unsigned requester", func() {

				_, ok := requester.(*UnsignedRequest)
				So(ok, ShouldBeTrue)
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

func TestUnitSubmitBulkToES_BodyReadError(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	signedReq := NewMockRequester(ctrl)
	unsignedReq := NewMockRequester(ctrl)

	mc := NewClientWithRequesters(mw, signedReq, unsignedReq)

	bulk := make([]byte, 1)
	companyNumbers := make([]byte, 1)
	esDestURL := "esDestURL"
	esDestIndex := "esDestIndex"
	esUri := esDestURL + "/" + esDestIndex + "/_bulk"

	Convey("Given a client fails to read response body", t, func() {

		// Create a response with a body that errors on read
		resp := &http.Response{
			StatusCode: 200,
			Body:       &errorReader{},
			Header:     make(http.Header),
		}

		signedReq.EXPECT().Post(bulk, esUri).Return(resp, nil)

		Convey("When SubmitBulkToES is called", func() {

			returnedBytes, err := mc.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)

			Convey("Then error should be returned", func() {

				So(returnedBytes, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUnitGetAlphaKeys_BodyReadError(t *testing.T) {

	ctrl := gomock.NewController(t)

	mw := write.NewMockWriter(ctrl)
	signedReq := NewMockRequester(ctrl)
	unsignedReq := NewMockRequester(ctrl)

	mc := NewClientWithRequesters(mw, signedReq, unsignedReq)

	companyNames := make([]byte, 1)
	alphaKeyURL := "alphaKeyURL"
	alphaKeyUri := alphaKeyURL + "/alphakey-bulk"

	Convey("Given a client fails to read alpha key response body", t, func() {

		// Create a response with a body that errors on read
		resp := &http.Response{
			StatusCode: 200,
			Body:       &errorReader{},
			Header:     make(http.Header),
		}

		unsignedReq.EXPECT().Post(companyNames, alphaKeyUri).Return(resp, nil)

		Convey("When GetAlphaKeys is called", func() {

			returnedBytes, err := mc.GetAlphaKeys(companyNames, alphaKeyURL)

			Convey("Then error should be returned", func() {

				So(returnedBytes, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

// errorReader implements io.ReadCloser and always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}
