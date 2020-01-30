package eshttp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/companieshouse/elasticsearch-data-loader/write"
)

const applicationJSON = "application/json"

// Client provides an interface with which to communicate with Elastic Search by way of HTTP requests
type Client interface {
	SubmitBulkToES(bulk []byte, companyNumbers []byte, esDestURL string, esDestIndex string) ([]byte, error)
	GetAlphaKeys(companyNames []byte, alphaKeyURL string) ([]byte, error)
}

// ClientImpl provides a concrete implementation of the Client interface
type ClientImpl struct {
	w write.Writer
	r Requester
}

// NewClient returns a concrete implementation of the Client interface
func NewClient(writer write.Writer) Client {

	return &ClientImpl{
		w: writer,
		r: NewRequester(),
	}
}

// NewClientWithRequester returns a concrete implementation of the Client interface, taking a custom Requester
func NewClientWithRequester(writer write.Writer, requester Requester) Client {

	return &ClientImpl{
		w: writer,
		r: requester,
	}
}

// SubmitBulkToES uses an HTTP post request to submit data to Elastic Search
func (c *ClientImpl) SubmitBulkToES(bulk []byte, companyNumbers []byte, esDestURL string, esDestIndex string) ([]byte, error) {

	uri := fmt.Sprintf("%s/%s/_bulk", esDestURL, esDestIndex)

	r, err := c.r.Post(bulk, uri)
	if err != nil {
		c.w.LogPostError(string(companyNumbers))
		log.Printf("error posting request %s: data %s", err, string(bulk))
		return nil, err
	}

	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body after posting bulk to ES: %s", err)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if r.StatusCode > 299 {
		c.w.LogUnexpectedResponse(string(companyNumbers))
		log.Printf("unexpected put response %s: data %s", r.Status, string(bulk))
		return nil, errors.New("invalid response")
	}

	return b, err
}

// GetAlphaKeys performs a POST request to fetch alpha keys for a given set of company names
func (c *ClientImpl) GetAlphaKeys(companyNames []byte, alphaKeyURL string) ([]byte, error) {

	uri := fmt.Sprintf("%s/alphakey-bulk", alphaKeyURL)

	r, err := c.r.Post(companyNames, uri)
	if err != nil {
		c.w.LogAlphaKeyErrors(string(companyNames))
		log.Printf("error fetching alpha keys %s: data %s", err, string(companyNames))
		return nil, err
	}

	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body after fetching alpha keys: %s", err)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, err
}
