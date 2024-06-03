package functions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// function handles valid POST request with JSON body and query parameter
func TestHandleDigestValidPostRequest(t *testing.T) {
	fmt.Println("!!! Starting digest summary calculation tests on digest.go!!!")

	// Create a new HTTP request with POST method
	reqBody := `{"digest": "SGVsbG8gd29ybGQ="}`
	req, err := http.NewRequest(http.MethodPost, "/digest/calculateSummary?hash=sha256", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleDigest function
	HandleDigest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := `{"digestSummary":"ZOyIygCyaOW6GjVnihtTFtIS9PNmskdyMlNKiuyjfzw=","URLSafeDigestSummary":"ZOyIygCyaOW6GjVnihtTFtIS9PNmskdyMlNKiuyjfzw=","algorithmUsed":"sha256"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, rr.Body.String())
	}
}

// function handles valid POST request with JSON body and default hash algorithm
func TestHandleDigestValidPostRequestDefaultHashAlgorithm(t *testing.T) {
	// Create a new HTTP request with POST method
	reqBody := `{"digest": "SGVsbG8gd29ybGQ="}`
	req, err := http.NewRequest(http.MethodPost, "/digest", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleDigest function
	HandleDigest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := `{"digestSummary":"ZOyIygCyaOW6GjVnihtTFtIS9PNmskdyMlNKiuyjfzw=","URLSafeDigestSummary":"ZOyIygCyaOW6GjVnihtTFtIS9PNmskdyMlNKiuyjfzw=","algorithmUsed":"sha256"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, rr.Body.String())
	}
}

// function handles valid POST request with JSON body and sha384
func TestHandleDigestValidPostRequestSha386(t *testing.T) {
	// Create a new HTTP request with POST method
	reqBody := `{"digest": "SGVsbG8gd29ybGQ="}`
	req, err := http.NewRequest(http.MethodPost, "/digest?hash=sha384", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleDigest function
	HandleDigest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := `{"digestSummary":"kgOwxEOf0eauWHiGYze3xTKs1tkmAVDIAxjoq4wnzjMBifjflPuJDfHSmP82Bifh","URLSafeDigestSummary":"kgOwxEOf0eauWHiGYze3xTKs1tkmAVDIAxjoq4wnzjMBifjflPuJDfHSmP82Bifh","algorithmUsed":"sha384"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, rr.Body.String())
	}
}

// function handles valid POST request with JSON body and sha512
func TestHandleDigestValidPostRequestUsingSha512(t *testing.T) {
	// Create a new HTTP request with POST method
	reqBody := `{"digest": "SGVsbG8gd29ybGQ="}`
	req, err := http.NewRequest(http.MethodPost, "/digest?hash=sha512", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleDigest function
	HandleDigest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := `{"digestSummary":"t/eDuu2Cl/DbkXRiGE/08I5pwtXl95qUJgD5cl9Yzh8pwYE5v4CwbA//K900c4RS7PQMSIwip+PYDN9vnBwNRw==","URLSafeDigestSummary":"t_eDuu2Cl_DbkXRiGE_08I5pwtXl95qUJgD5cl9Yzh8pwYE5v4CwbA__K900c4RS7PQMSIwip-PYDN9vnBwNRw==","algorithmUsed":"sha512"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, rr.Body.String())
	}
}

// function handles invalid request method
func TestHandleDigestInvalidRequestMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/digest/calculateSummary?hash=sha256", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	HandleDigest(rr, req)

	fmt.Printf("Actual status code: %d\n", rr.Code)            // Print actual status code
	fmt.Printf("Actual response body: %s\n", rr.Body.String()) // Print actual response body

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	expectedBody := "Invalid request method"
	actualBody := strings.TrimSpace(rr.Body.String())
	if actualBody != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, actualBody)
	}
}

// function handles invalid hash algorithm
func TestHandleDigestInvalidHashAlgorithm(t *testing.T) {
	// Create a new HTTP request with POST method and invalid hash algorithm
	reqBody := `{"digest": "SGVsbG8gd29ybGQ="}`
	req, err := http.NewRequest(http.MethodPost, "/digest?hash=md5", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleDigest function
	HandleDigest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}

	// Check the response body
	expectedBody := "Unsupported hash algorithm: md5"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, rr.Body.String())
	}
}

// function handles invalid base64 digest
func TestHandleDigestInvalidBase64Digest(t *testing.T) {
	// Create a new HTTP request with POST method
	reqBody := `{"digest": "invalid_base64_digest"}`
	req, err := http.NewRequest(http.MethodPost, "/digest", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the HandleDigest function
	HandleDigest(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}

	// Check the response body
	expectedBody := "Error decoding base64 digest: illegal base64 data at input byte 20"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %s, but got %s", expectedBody, rr.Body.String())
	}
}
