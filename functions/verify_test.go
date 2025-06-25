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

package functions

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// parseCertificate function tests:
func TestParseCertificateValidCertificate(t *testing.T) {
	println("!!! Starting digest verification tests on verify.go !!!")
	// This is your test certificate string. Replace it with the actual string.
	certificateStr := "MIIGRzCCBC+gAwIBAgIQbgfyMf1q2W9ck0+6CRSZSDANBgkqhkiG9w0BAQsFADCBhjELMAkGA1UEBhMCTFYxOTA3BgNVBAoMMFZBUyBMYXR2aWphcyBWYWxzdHMgcmFkaW8gdW4gdGVsZXbEq3ppamFzIGNlbnRyczEaMBgGA1UEYQwRTlRSTFYtNDAwMDMwMTEyMDMxIDAeBgNVBAMMF0RFTU8gZVBhcmFrc3RzIElDQSAyMDE3MB4XDTE5MDMyMTA4NDc1M1oXDTI0MDMyMTA4NDc1M1owXjELMAkGA1UEBhMCTFYxCzAJBgNVBAsMAklUMRAwDgYDVQQKDAdaWiBEYXRzMRowGAYDVQRhDBFOVFJMVi00MDAwMzI3ODQ2NzEUMBIGA1UEAwwLTW9iQXV0aENlcnQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDDLEbR2BDlHj521jYT3eiZD3JrbaVEv9iWqQG30xYyXTbU/sZzpf4o7yhSBcvyubsdbwVStS+LHzkoC17bAYvK4WjuHP33nkX8vcaL9lCRoxJZTCU/FBWOHnMkOY7Ifw1MUX7GbeOJCu0nDN7dix0s/w/PuNoa/mnFexubH2zaE9MoHO002DljskJTR2c/p3m9/jCveirCJYTW7hvICYkPNO494xHIt/uSq8M51fEJ27z9FPgy516Ea6veLb6vM52MBu5z8ZpIH7IW7kmdNWaWd+e61h6M6+YWjAe7UNNz8yUWZ4W1ViZHqKDqLB7TCjpAted/29YQIqG/8b4ZikpVAgMBAAGjggHWMIIB0jAMBgNVHRMBAf8EAjAAMA4GA1UdDwEB/wQEAwIHgDATBgNVHSUEDDAKBggrBgEFBQcDAjAdBgNVHQ4EFgQUrgiPlOna4q0FoI2EArZBE1vxoZkwHwYDVR0jBBgwFoAULZuJhBv1f9dJouY2Y3TZFBAO/FUwgYsGA1UdIASBgzCBgDA7BgYEAI96AQEwMTAvBggrBgEFBQcCARYjaHR0cHM6Ly93d3cuZXBhcmFrc3RzLmx2L3JlcG9zaXRvcnkwQQYMKwYBBAGB+j0CAgEBMDEwLwYIKwYBBQUHAgEWI2h0dHBzOi8vd3d3LmVwYXJha3N0cy5sdi9yZXBvc2l0b3J5MIGABggrBgEFBQcBAQR0MHIwRQYIKwYBBQUHMAKGOWh0dHA6Ly9kZW1vLmVwYXJha3N0cy5sdi9jZXJ0L2RlbW9fZVBhcmFrc3RzX0lDQV8yMDE3LmNydDApBggrBgEFBQcwAYYdaHR0cDovL29jc3AucHJlcC5lcGFyYWtzdHMubHYwTAYDVR0fBEUwQzBBoD+gPYY7aHR0cDovL2RlbW8uZXBhcmFrc3RzLmx2L2NybC9kZW1vX2VQYXJha3N0c19JQ0FfMjAxN18yMS5jcmwwDQYJKoZIhvcNAQELBQADggIBAB1QqSgtQw1j5Y7oi5I0PLuFoRwULnEeFBgosZxrAT7RW8mwH8RxFXn9eB0KaraA02Atlon8JuhsYEC69QbKuCDl57gzvCGqaf+JBigm/fymF55H6+1U0D8F4YJ/rS7quaz6X9+Oj0bANCE+M6jWKTTU4gIOrX125pfOHNOYUn9sVgxYzi7MTzLE2N0xBqlCmhTbTcGTWk0jQ39zhN+BpEUydeCEZmtJshUfIhFKCqglC47yHUw5KKrlSjEHx6MLguIph2OHJoqkJGxTJiIRFqWqEnTJJPJn45wEewjQW6E5HGq/52XghW2o2yszRC79yEGr4iSbAvx7h7Q8NidOE1IYBLVvE50E4++lgxXGRX5oKGaEp9tOye34wdQ2ubFO+wQMjpxFN2b6fYf7P+qYJZ77S+rxpcfRz9Wo4SYJ1nxopwpsNU/mRPv/UTZSVzkf1Qe909mZYhDFw/8GoGe+BJsZMUycwaYJb9gKE9ytaROQRjLDtKpige+Blb8kDlo6fsC4wZ1XqCXPscERzKE72g37Hw6FNdb4gRyA1q+ge2VYsXiEsy42NZ44oALuplLD38Gt4roEFMI3eELF2blOEJYDmUejWn8eAc0SSuA27urmkrHC4HBfzLzVEsH7TxKs9C3lnOZT+uUiDTbv6EaA3OFmlmv1ohzvdaO+LaKIqazQ"

	_, err := parseCertificate(certificateStr)

	// The test passes if err is nil (i.e., no error occurred, which means the certificate is valid).
	// If an error occurred (i.e., the certificate is not valid), the test fails.
	assert.NoError(t, err)
}

func TestParseCertificateEmptyCertificateString(t *testing.T) {
	certificateStr := ""

	certificate, err := parseCertificate(certificateStr)

	assert.Error(t, err)
	assert.Nil(t, certificate)
}

// decodeBase64 function tests
func TestDecodeBase64(t *testing.T) {
	// This is a test string.
	testString := "Hello, World!"

	// Encode the test string to base64.
	encodedString := base64.StdEncoding.EncodeToString([]byte(testString))

	// Call the function with the encoded string.
	decodedBytes, err := decodeBase64(encodedString)

	// The test passes if err is nil (i.e., no error occurred, which means the string was correctly decoded).
	// If an error occurred (i.e., the string was not correctly decoded), the test fails.
	assert.NoError(t, err)

	// The test also passes if the decoded bytes match the original string.
	// If they don't match, the test fails.
	assert.Equal(t, testString, string(decodedBytes))
}

func TestDecodeBase64Invalid(t *testing.T) {
	// This is an invalid base64 string.
	invalidString := "This is not a valid base64 string!"

	// Call the function with the invalid string.
	_, err := decodeBase64(invalidString)

	// The test passes if an error was returned (i.e., the string was not correctly decoded).
	// If no error was returned (i.e., the string was incorrectly decoded), the test fails.
	assert.Error(t, err)
}

// VerifyECDSASignature

func TestVerifyECDSASignature(t *testing.T) {
	// Generate a new key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Compute a digest of some data
	data := []byte("Hello, World!")
	hash := sha256.Sum256(data)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatal(err)
	}

	// Serialize the signature
	signature := append(r.Bytes(), s.Bytes()...)

	// Verify the signature
	err = verifyECDSASignature(&privateKey.PublicKey, hash[:], signature)

	// The test passes if err is nil (i.e., the signature is valid).
	// If an error occurred (i.e., the signature is not valid), the test fails.
	assert.NoError(t, err)
}

func TestVerifyECDSASignatureInvalid(t *testing.T) {
	// Generate a new key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Compute a digest of some data
	data := []byte("Hello, World!")
	hash := sha256.Sum256(data)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatal(err)
	}

	// Serialize the signature
	signature := append(r.Bytes(), s.Bytes()...)

	// Modify the signature to make it invalid
	signature[0] ^= 0xff

	// Verify the signature
	err = verifyECDSASignature(&privateKey.PublicKey, hash[:], signature)

	// The test passes if an error was returned (i.e., the signature is not valid).
	// If no error was returned (i.e., the signature was incorrectly verified), the test fails.
	assert.Error(t, err)
}

func TestVerifyECDSASignatureInvalidLength(t *testing.T) {
	// Generate a new key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Compute a digest of some data
	data := []byte("Hello, World!")
	hash := sha256.Sum256(data)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatal(err)
	}

	// Serialize the signature
	signature := append(r.Bytes(), s.Bytes()...)

	// Modify the signature to make its length invalid
	signature = append(signature, 0xff)

	// Verify the signature
	err = verifyECDSASignature(&privateKey.PublicKey, hash[:], signature)
	fmt.Println(err)
	// The test passes if an error was returned (i.e., the signature length is invalid).
	// If no error was returned (i.e., the signature length was incorrectly verified), the test fails.
	assert.Error(t, err)
}

func TestVerifyECDSASignatureFailedVerification(t *testing.T) {
	// Generate a new key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Compute a digest of some data
	data := []byte("Hello, World!")
	hash := sha256.Sum256(data)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatal(err)
	}

	// Serialize the signature
	signature := append(r.Bytes(), s.Bytes()...)

	// Modify the hash to make the signature verification fail
	hash[0] ^= 0xff

	// Verify the signature
	err = verifyECDSASignature(&privateKey.PublicKey, hash[:], signature)
	fmt.Println(err)
	// The test passes if an error was returned (i.e., the signature verification failed).
	// If no error was returned (i.e., the signature was incorrectly verified), the test fails.
	assert.Error(t, err)
}
