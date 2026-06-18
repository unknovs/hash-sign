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

package responses

// CSC API v2 response shapes (CSC API V2.2.0.0).

// CscInfoResponse is the POST /csc/v2/info response.
type CscInfoResponse struct {
	Specs                     string            `json:"specs"`
	Name                      string            `json:"name"`
	Region                    string            `json:"region,omitempty"`
	Lang                      string            `json:"lang,omitempty"`
	Description               string            `json:"description,omitempty"`
	AuthType                  []string          `json:"authType"`
	OAuth2                    string            `json:"oauth2,omitempty"`
	Methods                   []string          `json:"methods"`
	SignAlgorithms            CscSignAlgorithms `json:"signAlgorithms"`
	AsynchronousOperationMode bool              `json:"asynchronousOperationMode"`
}

type CscSignAlgorithms struct {
	Algos []string `json:"algos"`
}

// CscTokenResponse is the POST /csc/v2/oauth2/token response.
type CscTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// CscCredentialsListResponse is the POST /csc/v2/credentials/list response.
type CscCredentialsListResponse struct {
	CredentialIDs []string `json:"credentialIDs"`
}

// CscCredentialInfoResponse is the POST /csc/v2/credentials/info response.
type CscCredentialInfoResponse struct {
	Key                CscKey  `json:"key"`
	Cert               CscCert `json:"cert"`
	SignatureQualifier string  `json:"signatureQualifier,omitempty"`
	Auth               CscAuth `json:"auth"`
	SCAL               string  `json:"SCAL,omitempty"`
	Multisign          int     `json:"multisign"`
}

type CscKey struct {
	Status string   `json:"status"`
	Algo   []string `json:"algo"`
	Len    int      `json:"len,omitempty"`
}

type CscCert struct {
	Status       string   `json:"status"`
	Certificates []string `json:"certificates"` // base64-DER, leaf first
	IssuerDN     string   `json:"issuerDN,omitempty"`
	SerialNumber string   `json:"serialNumber,omitempty"`
	SubjectDN    string   `json:"subjectDN,omitempty"`
	ValidFrom    string   `json:"validFrom,omitempty"`
	ValidTo      string   `json:"validTo,omitempty"`
	QcStatements []string `json:"qcStatements,omitempty"`
}

type CscAuth struct {
	Mode string `json:"mode"`
}

// CscSignHashResponse is the POST /csc/v2/signatures/signHash response: the base64
// signatures in the same order as the request digests.
type CscSignHashResponse struct {
	Signatures []string `json:"signatures"`
}

// CscError is the CSC error body returned alongside a non-2xx HTTP status.
type CscError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}
