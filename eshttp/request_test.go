package eshttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

func TestUnitNewUnsignedRequester(t *testing.T) {
	// Test that NewUnsignedRequester always returns an unsigned requester
	requester := NewUnsignedRequester()

	if _, ok := requester.(*UnsignedRequest); !ok {
		t.Errorf("expected UnsignedRequest, got %T", requester)
	}

	// Verify the requester has an http client
	unsignedReq, _ := requester.(*UnsignedRequest)
	if unsignedReq.httpClient == nil {
		t.Error("expected httpClient to be initialized")
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

func TestUnitSignedRequest_PostWithValidSigner(t *testing.T) {
	// Create a test server that expects requests with Content-Type
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != applicationJSON {
			t.Errorf("expected Content-Type %s, got %s", applicationJSON, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a mock signer that records calls
	mockSigner := &mockSigner{}

	requester := &SignedRequest{
		signer:     mockSigner,
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

	if !mockSigner.signRequestCalled {
		t.Error("expected SignRequest to be called on signer")
	}
}

// mockSigner implements the Signer interface for testing
type mockSigner struct {
	signRequestCalled bool
	lastRequest       *http.Request
}

func (m *mockSigner) SignRequest(req *http.Request) error {
	m.signRequestCalled = true
	m.lastRequest = req
	// Add a marker header to show the request was signed
	req.Header.Set("X-Signed", "true")
	return nil
}

func TestUnitSignedRequest_PostErrorHandling_NilSigner(t *testing.T) {
	// Test that SignedRequest properly handles nil signer error
	requester := &SignedRequest{
		signer:     nil, // Intentionally nil to test error handling
		httpClient: &http.Client{},
	}

	body := []byte(`{"test": "data"}`)
	_, err := requester.Post(body, "http://localhost:9200")
	if err == nil || err.Error() != "AWS SigV4 signer not available" {
		t.Errorf("expected 'AWS SigV4 signer not available' error, got %v", err)
	}
}

func TestUnitSignedRequest_PostErrorHandling_NilHTTPClient(t *testing.T) {
	// Test that SignedRequest properly handles nil http client error
	mockSigner := &mockSigner{}

	requester := &SignedRequest{
		signer:     mockSigner,
		httpClient: nil,
	}

	body := []byte(`{"test": "data"}`)
	_, err := requester.Post(body, "http://localhost:9200")
	if err == nil || err.Error() != "HTTP client not available" {
		t.Errorf("expected 'HTTP client not available' error, got %v", err)
	}
}

func TestUnitSignedRequest_PostErrorHandling_InvalidURL(t *testing.T) {
	// Test error handling when URL is invalid
	mockSigner := &mockSigner{}

	requester := &SignedRequest{
		signer:     mockSigner,
		httpClient: &http.Client{},
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

	mockSigner := &mockSigner{}

	requester := &SignedRequest{
		signer:     mockSigner,
		httpClient: &http.Client{},
	}

	body := []byte(`{"test": "data"}`)
	resp, err := requester.Post(body, server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if mockSigner.lastRequest.Header.Get("Content-Type") != applicationJSON {
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
