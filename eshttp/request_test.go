package eshttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUnsignedRequest_Post(t *testing.T) {
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

func TestNewRequester_Unsigned(t *testing.T) {
	// Ensure USE_AWS_SIGV4 is not set
	os.Unsetenv("USE_AWS_SIGV4")

	requester := NewRequester()

	// Should return UnsignedRequest when env var is not set or false
	if _, ok := requester.(*UnsignedRequest); !ok {
		t.Errorf("expected UnsignedRequest, got %T", requester)
	}
}

func TestNewRequester_SignedNoCredentials(t *testing.T) {
	// Set env var to true but ensure no AWS credentials available
	t.Setenv("USE_AWS_SIGV4", "true")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	// Also unset AWS_PROFILE to prevent loading from ~/.aws/config
	t.Setenv("AWS_PROFILE", "")

	requester := NewRequester()

	// Should gracefully fall back to UnsignedRequest when credentials unavailable
	if _, ok := requester.(*UnsignedRequest); !ok {
		t.Errorf("expected graceful fallback to UnsignedRequest, got %T", requester)
	}
}

