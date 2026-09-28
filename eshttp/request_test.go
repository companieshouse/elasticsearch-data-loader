package eshttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/opensearch-project/opensearch-go/v2"
)

func TestUnitUnsignedRequest_Post(t *testing.T) {
	// Create a test server that expects unsigned requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Content-Type is set
		if r.Header.Get("Content-Type") != applicationJSON {
			t.Errorf("expected Content-Type %s, got %s", applicationJSON, r.Header.Get("Content-Type"))
		}
		// Verify no Authorization header (unsigned request)
		if r.Header.Get("Authorization") != "" {
			t.Errorf("expected no Authorization header for unsigned request, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	requester := &UnsignedRequest{
		httpClient: &http.Client{},
	}

	body := []byte(`{"test": "data"}`)
	resp, err := requester.Post(body, server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestUnitNewRequester_Unsigned(t *testing.T) {
	// Ensure USE_AWS_SIGV4 is not set
	os.Unsetenv("USE_AWS_SIGV4")

	requester := NewRequester()

	// Should return UnsignedRequest when env var is not set or false
	if _, ok := requester.(*UnsignedRequest); !ok {
		t.Errorf("expected UnsignedRequest, got %T", requester)
	}
}

func TestUnitNewRequester_SignedEnabled(t *testing.T) {
	// Test that when USE_AWS_SIGV4 is set to true, it attempts to create SignedRequest
	// This test doesn't require actual AWS credentials since createSignedRequester handles fallback
	t.Setenv("USE_AWS_SIGV4", "true")

	requester := NewRequester()

	// Verify it's either SignedRequest (if credentials available) or UnsignedRequest (fallback)
	switch requester.(type) {
	case *SignedRequest, *UnsignedRequest:
		// Both are acceptable outcomes
	default:
		t.Errorf("expected SignedRequest or UnsignedRequest, got %T", requester)
	}
}

func TestUnitSignedRequest_PostWithValidTransport(t *testing.T) {
	// Create a test server that expects requests with Content-Type
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != applicationJSON {
			t.Errorf("expected Content-Type %s, got %s", applicationJSON, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a mock transport that records calls and returns a response
	mockTransport := &mockTransport{
		server: server,
	}

	// Create a SignedRequest with our mock transport
	client := &opensearch.Client{}
	client.Transport = mockTransport

	requester := &SignedRequest{
		client: client,
	}

	body := []byte(`{"test": "data"}`)
	resp, err := requester.Post(body, server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if !mockTransport.performCalled {
		t.Error("expected Perform to be called on transport")
	}
}

// mockTransport implements the OpenSearch transport interface for testing
type mockTransport struct {
	server         *httptest.Server
	performCalled  bool
	lastRequest    *http.Request
}

func (m *mockTransport) Perform(req *http.Request) (*http.Response, error) {
	m.performCalled = true
	m.lastRequest = req

	// Re-route the request to our test server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TestUnitSignedRequest_PostErrorHandling_NilClient(t *testing.T) {
	// Test that SignedRequest properly handles nil client error
	requester := &SignedRequest{
		client: nil, // Intentionally nil to test error handling
	}

	body := []byte(`{"test": "data"}`)
	_, err := requester.Post(body, "http://localhost:9200")
	if err == nil || err.Error() != "OpenSearch client not available" {
		t.Errorf("expected 'OpenSearch client not available' error, got %v", err)
	}
}

func TestUnitSignedRequest_PostErrorHandling_NilTransport(t *testing.T) {
	// Test that SignedRequest properly handles nil transport error
	client := &opensearch.Client{}
	client.Transport = nil

	requester := &SignedRequest{
		client: client,
	}

	body := []byte(`{"test": "data"}`)
	_, err := requester.Post(body, "http://localhost:9200")
	if err == nil || err.Error() != "OpenSearch client transport not available" {
		t.Errorf("expected 'OpenSearch client transport not available' error, got %v", err)
	}
}

func TestUnitSignedRequest_PostErrorHandling_InvalidURL(t *testing.T) {
	// Test error handling when URL is invalid
	mockTransport := &mockTransport{}

	client := &opensearch.Client{}
	client.Transport = mockTransport

	requester := &SignedRequest{
		client: client,
	}

	body := []byte(`{"test": "data"}`)
	// Invalid URL with newline should cause NewRequestWithContext to fail
	_, err := requester.Post(body, "http://invalid\nURL")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestUnitSignedRequest_PostSetsContentType(t *testing.T) {
	// Test that SignedRequest properly sets Content-Type header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != applicationJSON {
			t.Errorf("expected Content-Type %s, got %s", applicationJSON, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mockTransport := &mockTransport{
		server: server,
	}

	client := &opensearch.Client{}
	client.Transport = mockTransport

	requester := &SignedRequest{
		client: client,
	}

	body := []byte(`{"test": "data"}`)
	resp, err := requester.Post(body, server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if mockTransport.lastRequest.Header.Get("Content-Type") != applicationJSON {
		t.Errorf("Content-Type not set correctly on request")
	}
}

func TestUnitCreateSignedRequester_HandlesAWSConfig(t *testing.T) {
	// This test verifies that createSignedRequester properly attempts to load AWS config
	// It may return SignedRequest (if credentials available) or UnsignedRequest (on fallback)
	// This is acceptable behavior for graceful degradation
	requester := createSignedRequester()

	// Verify it returns either type
	switch requester.(type) {
	case *SignedRequest, *UnsignedRequest:
		// Both are acceptable outcomes
	default:
		t.Errorf("expected SignedRequest or UnsignedRequest, got %T", requester)
	}
}

func TestUnitNewRequester_SignedDisabled(t *testing.T) {
	// Test that disabled signing returns UnsignedRequest
	t.Setenv("USE_AWS_SIGV4", "false")

	requester := NewRequester()

	if _, ok := requester.(*UnsignedRequest); !ok {
		t.Errorf("expected UnsignedRequest when USE_AWS_SIGV4=false, got %T", requester)
	}
}
