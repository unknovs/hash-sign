package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/routes/responses"
)

const endpoint = "/certificates?key=rsa&type=auth"
const expectedStatusWrong = "Expected status code %d, but got %d"
const failedRequest = "Failed to create request: %v"

func TestHandleCertificatesRequestGetRequestWithCertificates(t *testing.T) {

	fmt.Println("!!! Starting certificate handler tests on digest.go!!!")
	env.RsaAuthCert = "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"
	t.Cleanup(func() {
		env.RsaAuthCert = ""
	})

	// Setup
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf(failedRequest, err)
	}
	rr := httptest.NewRecorder()

	// Execute
	HandleCertificatesRequest(rr, req)

	// Verify
	if rr.Code != http.StatusOK {
		t.Errorf(expectedStatusWrong, http.StatusOK, rr.Code)
	}
	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type header %s, but got %s", expectedContentType, rr.Header().Get("Content-Type"))
	}
	expectedResponse := responses.CertificatesResponse{
		RsaAuthenticationCertificate: env.RsaAuthCert,
	}
	var actualResponse responses.CertificatesResponse
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		t.Errorf("Expected response %+v, but got %+v", expectedResponse, actualResponse)
	}
}

func TestHandleCertificatesRequestNonGetRequest(t *testing.T) {
	// Setup
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		t.Fatalf(failedRequest, err)
	}
	rr := httptest.NewRecorder()

	// Execute
	HandleCertificatesRequest(rr, req)

	// Verify
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf(expectedStatusWrong, http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleCertificatesRequestNoCertificatesExist(t *testing.T) {
	// Setup
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf(failedRequest, err)
	}
	rr := httptest.NewRecorder()
	env.RsaAuthCert = ""
	env.RsaSigningCert = ""
	env.EcdsaAuthCert = ""
	env.EcdsaSigningCert = ""

	// Execute
	HandleCertificatesRequest(rr, req)

	// Verify
	if rr.Code != http.StatusNotFound {
		t.Errorf(expectedStatusWrong, http.StatusNotFound, rr.Code)
	}
	expectedErrorMessage := "No certificates found in environment"
	var actualResponse map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if actualResponse["error"] != expectedErrorMessage {
		t.Errorf("Expected error message %s, but got %s", expectedErrorMessage, actualResponse["error"])
	}
}
