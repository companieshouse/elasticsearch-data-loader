package eshttp

import (
	"bytes"
	"net/http"
)

// Requester provides an interface by which to execute HTTP requests
type Requester interface {
	PostBulkToElasticSearch(bulk []byte, esDestURL string, esDestIndex string) (*http.Response, error)
}

// Request provides a concrete implementation of the Requester interface
type Request struct{}

// NewRequester returns a concrete implementation of the Requester interface
func NewRequester() Requester {

	return &Request{}
}

// PostBulkToElasticSearch posts a 'bulk' to Elastic Search
func (req *Request) PostBulkToElasticSearch(bulk []byte, esDestURL string, esDestIndex string) (*http.Response, error) {

	return http.Post(esDestURL+"/"+esDestIndex+"/_bulk", applicationJSON, bytes.NewReader(bulk))
}
