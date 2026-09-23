package eshttp

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
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
	httpClient   *http.Client
	credentials  aws.Credentials
	signer       *v4.Signer
	region       string
}

// NewRequester returns a concrete implementation of the Requester interface
func NewRequester() Requester {

	useAWSSignV4 := os.Getenv("USE_AWS_SIGV4") == "true"
	var credentials aws.Credentials
	var signer *v4.Signer
	region := "eu-west-2" // default region

	// Only load AWS config if SigV4 is enabled
	if useAWSSignV4 {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Printf("warning: failed to load AWS config for SigV4: %v, falling back to unsigned requests", err)
			useAWSSignV4 = false
		} else {
			// Get the credentials from the config
			creds, err := cfg.Credentials.Retrieve(context.Background())
			if err != nil {
				log.Printf("warning: failed to retrieve AWS credentials: %v, falling back to unsigned requests", err)
				useAWSSignV4 = false
			} else {
				credentials = creds
				signer = v4.NewSigner()
				if cfg.Region != "" {
					region = cfg.Region
				}
			}
		}
	}

	return &Request{
		useAWSSignV4: useAWSSignV4,
		httpClient:   &http.Client{},
		credentials:  credentials,
		signer:       signer,
		region:       region,
	}
}

// Post performs a POST request, using a provided body, against a given uri
// If USE_AWS_SIGV4 env var is set to "true", signs request with AWS SigV4
func (req *Request) Post(body []byte, uri string) (*http.Response, error) {

	if !req.useAWSSignV4 {
		// Standard unsigned request for local/non-AWS environments
		return http.Post(uri, applicationJSON, bytes.NewReader(body))
	}

	// Create request for signing
	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", applicationJSON)

	// Sign the request with AWS SigV4
	// Service name is "es" for Elasticsearch/OpenSearch (Amazon OpenSearch Service)
	err = req.signer.SignHTTP(context.Background(), req.credentials, httpReq, OpenSearchSigningServiceName, req.region, nil)
	if err != nil {
		log.Printf("error signing request with SigV4: %v", err)
		return nil, err
	}

	// Execute signed request
	return req.httpClient.Do(httpReq)
}

