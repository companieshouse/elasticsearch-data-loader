package eshttp

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

// OpenSearchSigningServiceName is the IAM action/service name prefix used by Amazon OpenSearch Service for SigV4 signing
// Used in requests like es:ESHttpPost
const OpenSearchSigningServiceName = "es"

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
	signer     *requestsigner.Signer
	httpClient *http.Client
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

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("warning: failed to load AWS config for SigV4: %v, falling back to unsigned requests", err)
		return &UnsignedRequest{
			httpClient: &http.Client{},
		}
	}

	signer, err := requestsigner.NewSigner(cfg)
	if err != nil {
		log.Printf("warning: failed to create AWS SigV4 signer: %v, falling back to unsigned requests", err)
		return &UnsignedRequest{
			httpClient: &http.Client{},
		}
	}

	log.Printf("info: AWS SigV4 signing enabled for OpenSearch requests")
	return &SignedRequest{
		signer:     signer,
		httpClient: &http.Client{},
	}
}

// Post performs an unsigned POST request
func (req *UnsignedRequest) Post(body []byte, uri string) (*http.Response, error) {
	return http.Post(uri, applicationJSON, bytes.NewReader(body))
}

// Post performs a SigV4-signed POST request for OpenSearch
func (req *SignedRequest) Post(body []byte, uri string) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", applicationJSON)

	// Use the OpenSearch client directly with the signer
	// This leverages the official OpenSearch Go client's signing mechanism
	err = req.signer.SignRequest(context.Background(), httpReq)
	if err != nil {
		log.Printf("error signing request with SigV4: %v", err)
		return nil, err
	}

	return req.httpClient.Do(httpReq)
}


