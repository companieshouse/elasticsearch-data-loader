package eshttp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/opensearch-project/opensearch-go/v2"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

// OpenSearchSigningServiceName is the IAM action/service name prefix used by Amazon OpenSearch Service for SigV4 signing
// Used in requests like es:ESHttpPost
const OpenSearchSigningServiceName = "es"

// Requester provides an interface by which to execute HTTP requests
type Requester interface {
	Post(body []byte, uri string) (*http.Response, error)
}

// Request provides a concrete implementation of the Requester interface
type Request struct {
	useAWSSignV4 bool
	client       *opensearch.Client
	httpClient   *http.Client
}

// NewRequester returns a concrete implementation of the Requester interface
func NewRequester() Requester {

	useAWSSignV4 := os.Getenv("USE_AWS_SIGV4") == "true"
	var osClient *opensearch.Client
	var httpClient *http.Client = &http.Client{}

	// If SigV4 is enabled, create OpenSearch client with AWS signer
	if useAWSSignV4 {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Printf("warning: failed to load AWS config for SigV4: %v, falling back to unsigned requests", err)
			useAWSSignV4 = false
		} else {
			signer, err := requestsigner.NewSigner(cfg)
			if err != nil {
				log.Printf("warning: failed to create AWS SigV4 signer: %v, falling back to unsigned requests", err)
				useAWSSignV4 = false
			} else {
				// Extract endpoint from environment or use placeholder
				// The actual endpoint URL will be provided per-request in the Post() method
				osClient, err = opensearch.NewClient(opensearch.Config{
					Addresses: []string{"https://placeholder"}, // Will be overridden per-request
					Signer:    signer,
				})
				if err != nil {
					log.Printf("warning: failed to create OpenSearch client: %v, falling back to unsigned requests", err)
					useAWSSignV4 = false
				} else {
					log.Printf("info: AWS SigV4 signing enabled for OpenSearch requests")
				}
			}
		}
	}

	return &Request{
		useAWSSignV4: useAWSSignV4,
		client:       osClient,
		httpClient:   httpClient,
	}
}

// Post performs a POST request, using a provided body, against a given uri
// If USE_AWS_SIGV4 env var is set to "true", signs request with AWS SigV4 via OpenSearch client
func (req *Request) Post(body []byte, uri string) (*http.Response, error) {

	if !req.useAWSSignV4 {
		// Standard unsigned request for local/non-AWS environments
		return http.Post(uri, applicationJSON, bytes.NewReader(body))
	}

	// Use OpenSearch client's signing transport for SigV4-signed requests
	// Create a new request that will be signed by the OpenSearch client's transport
	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", applicationJSON)

	// Use the OpenSearch client's HTTP transport which applies SigV4 signing
	if req.client != nil && req.client.Transport != nil {
		return req.client.Transport.RoundTrip(httpReq)
	}

	// Fallback to unsigned request if transport not available
	log.Printf("warning: OpenSearch client transport not available, falling back to unsigned request")
	return req.httpClient.Do(httpReq)
}


