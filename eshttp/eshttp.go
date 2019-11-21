package eshttp

import (
	"bytes"
	"errors"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"io/ioutil"
	"log"
	"net/http"
)

const applicationJSON = "application/json"

// Client provides an interface with which to communicate with Elastic Search by way of HTTP requests
type Client interface {
	SubmitDataToES(bulk []byte, bunchOfNamesAndNumbers []byte, esDestURL string, esDestIndex string) ([]byte, error)
}

// ClientImpl provides a concrete implementation of the Client interface
type ClientImpl struct {
	w write.Writer
}

// NewClient returns a concrete implementation of the Client interface
func NewClient(writer write.Writer) Client {

	return &ClientImpl{
		w: writer,
	}
}

// SubmitDataToES uses an HTTP post request to submit data to Elastic Search
func (c *ClientImpl) SubmitDataToES(bulk []byte, bunchOfNamesAndNumbers []byte, esDestURL string, esDestIndex string) ([]byte, error) {

	r, err := http.Post(esDestURL+"/"+esDestIndex+"/_bulk", applicationJSON, bytes.NewReader(bulk))
	if err != nil {
		c.w.LogPostError(string(bunchOfNamesAndNumbers))
		log.Printf("error posting request %s: data %s", err, string(bulk))
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error reading response body: %s", err)
	}

	if r.StatusCode > 299 {
		c.w.LogUnexpectedResponse(string(bunchOfNamesAndNumbers))
		log.Printf("unexpected put response %s: data %s", r.Status, string(bulk))
		return nil, errors.New("invalid response")
	}

	return b, err
}
