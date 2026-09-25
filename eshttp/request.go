package eshttp

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/opensearch-project/opensearch-go/v2"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

// Requester provides an interface by which to execute HTTP requests
type Requester interface {
	Post(body []byte, uri string) (*http.Response, error)
}

// UnsignedRequest provides unsigned HTTP requests (used for alpha-key and other non-OpenSearch services)
type UnsignedRequest struct {
	httpClient *http.Client
}

// SignedRequest provides AWS SigV4-signed HTTP requests (used for OpenSearch only)
type SignedRequest struct {
	client *opensearch.Client
}

// NewRequester returns the appropriate Requester based on USE_AWS_SIGV4 env var
// Returns UnsignedRequest for local/non-AWS environments
// Returns SignedRequest for AWS OpenSearch environments
func NewRequester() Requester {
	useAWSSignV4 := os.Getenv("USE_AWS_SIGV4") == "true"

	if !useAWSSignV4 {
		return &UnsignedRequest{
			httpClient: &http.Client{},
		}
	}

	return createSignedRequester()
}

// createSignedRequester creates a SignedRequest or falls back to UnsignedRequest on error
func createSignedRequester() Requester {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("warning: failed to load AWS config for SigV4: %v, falling back to unsigned requests", err)
		return &UnsignedRequest{httpClient: &http.Client{}}
	}

	signer, err := requestsigner.NewSigner(cfg)
	if err != nil {
		log.Printf("warning: failed to create AWS SigV4 signer: %v, falling back to unsigned requests", err)
		return &UnsignedRequest{httpClient: &http.Client{}}
	}

	// Note: client is only used to access the Transport for signing.
	// The actual request execution still uses our custom Post method.
	client, err := opensearch.NewClient(opensearch.Config{
		Signer: signer,
	})
	if err != nil {
		log.Printf("warning: failed to create OpenSearch client: %v, falling back to unsigned requests", err)
		return &UnsignedRequest{httpClient: &http.Client{}}
	}

	log.Printf("info: AWS SigV4 signing enabled for OpenSearch requests")
	return &SignedRequest{client: client}
}

// Post performs an unsigned POST request
func (req *UnsignedRequest) Post(body []byte, uri string) (*http.Response, error) {
	return http.Post(uri, applicationJSON, bytes.NewReader(body))
}

// Post performs a SigV4-signed POST request for OpenSearch
func (req *SignedRequest) Post(body []byte, uri string) (*http.Response, error) {
	if req.client == nil {
		return nil, errors.New("OpenSearch client not available")
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", applicationJSON)

	// Use the OpenSearch client's transport to sign and execute the request
	transport := req.client.Transport
	if transport == nil {
		return nil, errors.New("OpenSearch client transport not available")
	}

	return transport.Perform(httpReq)
}
