package functions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknovs/hash-sign/env"
)

// testing isGetMethod method
func TestIsGetMethodRequestMethodIsGetReturnsTrue(t *testing.T) {
	fmt.Println("!!! Starting global logic tests on logic_global.go!!!")
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	result := isGetMethod(request)
	if !result {
		t.Errorf("Expected true, but got false")
	}
}

func TestIsGetMethodRequestMethodIsNotGetReturnsFalse(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/", nil)
	result := isGetMethod(request)
	if result {
		t.Errorf("Expected false, but got true")
	}
}

func TestIsGetMethodInvalidRequestMethodReturnsFalse(t *testing.T) {
	request, _ := http.NewRequest("INVALID", "/", nil)
	result := isGetMethod(request)
	if result {
		t.Errorf("Expected false, but got true")
	}
}

// func TestIsGetMethod_LeadingTrailingWhiteSpacesInMethod_ReturnsTrue(t *testing.T) {
// 	request, _ := http.NewRequest("  GET  ", "/", nil)
// 	result := isGetMethod(request)
// 	if !result {
// 		t.Errorf("Expected true, but got false")
// 	}
// }

// testing isPostMethod method
func TestIsPostMethod(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected bool
	}{
		{
			name:     "Request Method Is Post",
			method:   http.MethodPost,
			expected: true,
		},
		{
			name:     "using GET instead of POST in POST method",
			method:   http.MethodGet,
			expected: false,
		},
		{
			name:     "Using Invalid Request Method",
			method:   "INVALID",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}
			result := isPostMethod(request)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

// testing apikey
func TestValidAPIKey(t *testing.T) {
	env.ApiKey = "test-api-key"
	t.Cleanup(func() {
		env.ApiKey = ""
	})

	// Create a new request with a valid API key in the header
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("API-Key", "test-api-key")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a handler function that will be called by APIKeyAuthorization
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Call the APIKeyAuthorization function with the handler function
	APIKeyAuthorization(handler).ServeHTTP(rr, req)

	// Check if the response status code is 200
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
	fmt.Println("Below API-Key from env variables:")
	fmt.Println(env.ApiKey)
	fmt.Printf("status code: %d\n", rr.Code)
	fmt.Printf("response body: %q\n", rr.Body)
	// Check if the response body is "OK"
	expectedBody := "OK"
	actualBody := rr.Body.String()
	if actualBody != expectedBody {
		t.Errorf("Expected response body %q, but got %q", expectedBody, actualBody)
	}
}

func TestEmptyAPIKeyV2(t *testing.T) {
	env.ApiKey = "test-api-key"
	t.Cleanup(func() {
		env.ApiKey = ""
	})

	// Create a new request with an empty API key in the header
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("API-Key", "")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a handler function that will be called by APIKeyAuthorization
	handler := http.HandlerFunc(APIKeyAuthorization(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Call the APIKeyAuthorization function with the handler function
	APIKeyAuthorization(handler).ServeHTTP(rr, req)

	// Check if the response status code is 401
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusUnauthorized, rr.Code)
	}
	fmt.Println("Below API-Key from env variables:")
	fmt.Println(env.ApiKey)
	fmt.Printf("status code: %d\n", rr.Code)
	fmt.Printf("response body: %q\n", rr.Body)
	// Check if the response body is "Unauthorized"
	expectedBody := "Unauthorized\n"
	actualBody := rr.Body.String()
	if actualBody != expectedBody {
		t.Errorf("Expected response body %q, but got %q", expectedBody, actualBody)
	}
}
