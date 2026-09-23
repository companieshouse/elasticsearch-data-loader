package eshttp

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"reflect"

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

// Request provides a concrete implementation of the Requester interface
type Request struct {
	useAWSSignV4 bool
	signer       interface{} // *requestsigner.Signer - type not exported, so use interface{}
	httpClient   *http.Client
}

// NewRequester returns a concrete implementation of the Requester interface
func NewRequester() Requester {

	useAWSSignV4 := os.Getenv("USE_AWS_SIGV4") == "true"
	var signer interface{}
	var httpClient *http.Client = &http.Client{}

	// If SigV4 is enabled, create AWS signer
	if useAWSSignV4 {
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Printf("warning: failed to load AWS config for SigV4: %v, falling back to unsigned requests", err)
			useAWSSignV4 = false
		} else {
			var signErr error
			signer, signErr = requestsigner.NewSigner(cfg)
			if signErr != nil {
				log.Printf("warning: failed to create AWS SigV4 signer: %v, falling back to unsigned requests", signErr)
				useAWSSignV4 = false
			} else {
				log.Printf("info: AWS SigV4 signing enabled for OpenSearch requests")
			}
		}
	}

	return &Request{
		useAWSSignV4: useAWSSignV4,
		signer:       signer,
		httpClient:   httpClient,
	}
}

// Post performs a POST request, using a provided body, against a given uri
// If USE_AWS_SIGV4 env var is set to "true", signs request with AWS SigV4
func (req *Request) Post(body []byte, uri string) (*http.Response, error) {

	if !req.useAWSSignV4 {
		// Standard unsigned request for local/non-AWS environments
		return http.Post(uri, applicationJSON, bytes.NewReader(body))
	}

	// Create request to be signed
	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", applicationJSON)

	// Sign the request with AWS SigV4 using OpenSearch signer
	if req.signer != nil {
		// Use reflection to call SignHTTP since Signer type isn't exported
		signerValue := reflect.ValueOf(req.signer)
		method := signerValue.MethodByName("SignHTTP")
		if method.IsValid() {
			// Call signer.SignHTTP(ctx, httpReq)
			results := method.Call([]reflect.Value{
				reflect.ValueOf(context.Background()),
				reflect.ValueOf(httpReq),
			})
			if len(results) > 0 {
				if errInterface := results[0].Interface(); errInterface != nil {
					if err, ok := errInterface.(error); ok && err != nil {
						log.Printf("warning: failed to sign request with SigV4: %v, sending unsigned request", err)
						return req.httpClient.Do(httpReq)
					}
				}
			}
		} else {
			log.Printf("warning: signer.SignHTTP method not found, sending unsigned request")
			return req.httpClient.Do(httpReq)
		}
	}

	// Execute signed request
	return req.httpClient.Do(httpReq)
}


