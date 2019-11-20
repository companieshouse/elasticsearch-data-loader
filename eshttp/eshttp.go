package eshttp

import (
	"bytes"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	esDestIndex     = "companies"
	esDestURL       = "http://localhost:9200"
	applicationJSON = "application/json"
)

// Client provides an interface with which to make HTTP calls to Elastic Search
type Client interface {
	SubmitToES(bulk []byte, bunchOfNamesAndNumbers []byte) ([]byte, error)
}

// ClientImpl implements the Client interface
type ClientImpl struct {
	writer write.Writer
}

// NewClient returns a concrete implementation of the Client interface
func NewClient(w write.Writer) Client {

	return &ClientImpl{
		writer: w,
	}
}

// SubmitToES submits data to Elastic Search by use of HTTP calls
func (c ClientImpl) SubmitToES(bulk []byte, bunchOfNamesAndNumbers []byte) ([]byte, error) {

	r, err := http.Post(esDestURL+"/"+esDestIndex+"/_bulk", applicationJSON, bytes.NewReader(bulk))
	if err != nil {
		c.writer.WriteToFile1(string(bunchOfNamesAndNumbers))
		log.Printf("error posting request %s: data %s", err, string(bulk))
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error reading response body: %s", err)
	}

	if r.StatusCode > 299 {
		c.writer.WriteToFile2(string(bunchOfNamesAndNumbers))
		log.Printf("unexpected put response %s: data %s", r.Status, string(bulk))
		return nil, err
	}

	return b, err
}
