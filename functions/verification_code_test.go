// SPDX-License-Identifier: MIT

// Copyright (c) 2024 Gatis Beikerts
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package functions_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unknovs/hash-sign/functions"
	"github.com/unknovs/hash-sign/routes/requests"
	"github.com/unknovs/hash-sign/routes/responses"
)

func TestCalculateVerificationCode(t *testing.T) {
	fmt.Println("!!! Starting verification code generation tests on verification_code.go !!!")
	// Prepare a request
	reqBody := requests.RequestVerificationCode{
		Hash: base64.StdEncoding.EncodeToString([]byte("test")),
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/digest/verificationCode", bytes.NewBuffer(reqBodyBytes))

	// Use httptest to record the response
	recorder := httptest.NewRecorder()

	// Call the function
	functions.CalculateVerificationCode(recorder, req)

	// Check the status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the response
	var resBody responses.VerificationCodeResponse
	_ = json.Unmarshal(recorder.Body.Bytes(), &resBody)

	// Add more checks here based on your requirements
	// For example, check the VerificationCode is in the expected range
	if resBody.VerificationCode < 0 || resBody.VerificationCode >= 10000 {
		t.Errorf("unexpected verification code: got %v", resBody.VerificationCode)
	}
}

func TestCalculateVerificationCodeForInvalidRequestMethod(t *testing.T) {
	// Test case: Invalid request method
	req, _ := http.NewRequest("GET", "/digest/verificationCode", nil)
	recorder := httptest.NewRecorder()
	functions.CalculateVerificationCode(recorder, req)
	if status := recorder.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

}

func TestCalculateVerificationCodeErrorForReadingRequestBody(t *testing.T) {
	// Test case: Error reading request body
	req, _ := http.NewRequest("POST", "/digest/verificationCode", bytes.NewBuffer([]byte("bad request body")))
	recorder := httptest.NewRecorder()
	functions.CalculateVerificationCode(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCalculateVerificationCodeErrorUnmarshalingRequestBody(t *testing.T) {
	// Test case: Error unmarshalling request body
	reqBody := "this is not a valid request body"
	req, _ := http.NewRequest("POST", "/digest/verificationCode", bytes.NewBuffer([]byte(reqBody)))
	recorder := httptest.NewRecorder()
	functions.CalculateVerificationCode(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCalculateVerificationCodeErrorDecodingBase64(t *testing.T) {
	// Test case: Error decoding base64 hash
	reqBodyBytes, _ := json.Marshal(requests.RequestVerificationCode{Hash: "invalid base64"})
	req, _ := http.NewRequest("POST", "/digest/verificationCode", bytes.NewBuffer(reqBodyBytes))
	recorder := httptest.NewRecorder()
	functions.CalculateVerificationCode(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

}

func TestCalculateVerificationCodeErrorReadingRequestBody(t *testing.T) {
	// Test case: Error reading request body
	req, _ := http.NewRequest("POST", "/digest/verificationCode", bytes.NewBuffer([]byte("bad request body")))
	recorder := httptest.NewRecorder()
	functions.CalculateVerificationCode(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

}
