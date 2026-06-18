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

// A minimal but spec-faithful CSC API v2 (CSC API V2.2.0.0) test server, exposed
// only when CSC_MOCK_ENABLED. It implements the request/response shapes of the
// methods a remote-signing client uses — info, oauth2 authorize/token,
// credentials list/info, signatures/signHash — so any CSC client can drive it.
// signHash is backed by this service's ECC key and credentials/info by
// ECDSA_SIGN_CERT.
//
// Authentication is not enforced (oauth2/authorize auto-approves and any bearer
// is accepted): it is a test mock, so the routes are registered without the
// API-Key wrapper.

package functions

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/routes/requests"
	"github.com/unknovs/hash-sign/routes/responses"
)

const (
	cscSpecVersion      = "2.2.0.0"
	cscMockCredentialID = "ecc-sign-1"
	cscECPublicKeyOID   = "1.2.840.10045.2.1" // id-ecPublicKey
	cscAlgoECDSASHA256  = "1.2.840.10045.4.3.2"
	cscAlgoECDSASHA384  = "1.2.840.10045.4.3.3"
	cscAlgoECDSASHA512  = "1.2.840.10045.4.3.4"
	cscMockAuthCode     = "csc-mock-code"
	cscMockAccessToken  = "csc-mock-access-token"
	cscDefaultQualifier = "eu_eidas_qes"
)

// oidQCStatements is id-pe-qcStatements (RFC 3739).
var oidQCStatements = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 1, 3}

// CscInfoHandler — POST /csc/v2/info (the one mandatory CSC method).
func CscInfoHandler(w http.ResponseWriter, r *http.Request) {
	if !isPostMethod(r) {
		cscError(w, http.StatusMethodNotAllowed, "invalid_request", "use POST")
		return
	}

	cscWriteJSON(w, http.StatusOK, responses.CscInfoResponse{
		Specs:                     cscSpecVersion,
		Name:                      "hash-sign CSC mock",
		AuthType:                  []string{"oauth2code"},
		Methods:                   []string{"credentials/list", "credentials/info", "signatures/signHash"},
		SignAlgorithms:            responses.CscSignAlgorithms{Algos: []string{cscAlgoECDSASHA256, cscAlgoECDSASHA384, cscAlgoECDSASHA512}},
		AsynchronousOperationMode: false,
	})
}

// CscAuthorizeHandler — GET /csc/v2/oauth2/authorize. Auto-approves and
// 302-redirects to redirect_uri?code=...&state=... (PKCE/client not enforced).
func CscAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	if !isGetMethod(r) {
		cscError(w, http.StatusMethodNotAllowed, "invalid_request", "use GET")
		return
	}

	q := r.URL.Query()
	redirectURI := q.Get("redirect_uri")
	if redirectURI == "" {
		cscError(w, http.StatusBadRequest, "invalid_request", "missing redirect_uri")
		return
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		cscError(w, http.StatusBadRequest, "invalid_request", "invalid redirect_uri")
		return
	}

	rq := u.Query()
	rq.Set("code", cscMockAuthCode)
	if state := q.Get("state"); state != "" {
		rq.Set("state", state)
	}
	u.RawQuery = rq.Encode()

	log.Printf("csc-mock: authorize -> 302 %s", u.String())
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// CscTokenHandler — POST /csc/v2/oauth2/token. Returns an opaque bearer token
// (code/PKCE/client credentials are not validated).
func CscTokenHandler(w http.ResponseWriter, r *http.Request) {
	if !isPostMethod(r) {
		cscError(w, http.StatusMethodNotAllowed, "invalid_request", "use POST")
		return
	}

	cscWriteJSON(w, http.StatusOK, responses.CscTokenResponse{
		AccessToken: cscMockAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})
}

// CscCredentialsListHandler — POST /csc/v2/credentials/list. One synthetic
// credential (this service holds a single ECC signing key/cert).
func CscCredentialsListHandler(w http.ResponseWriter, r *http.Request) {
	if !isPostMethod(r) {
		cscError(w, http.StatusMethodNotAllowed, "invalid_request", "use POST")
		return
	}

	cscWriteJSON(w, http.StatusOK, responses.CscCredentialsListResponse{CredentialIDs: []string{cscMockCredentialID}})
}

// CscCredentialsInfoHandler — POST /csc/v2/credentials/info. Returns the signing
// certificate (base64-DER) in cert.certificates and the key/cert details read
// from it (subject, validity, qcStatements, key length) — all derived from
// ECDSA_SIGN_CERT.
func CscCredentialsInfoHandler(w http.ResponseWriter, r *http.Request) {
	if !isPostMethod(r) {
		cscError(w, http.StatusMethodNotAllowed, "invalid_request", "use POST")
		return
	}

	if env.EcdsaSigningCert == "" {
		cscError(w, http.StatusNotFound, "invalid_request", "no signing certificate configured")
		return
	}

	der, err := cscCertToDER(env.EcdsaSigningCert)
	if err != nil {
		cscError(w, http.StatusInternalServerError, "server_error", "signing certificate is not valid base64/PEM")
		return
	}

	cert := responses.CscCert{
		Status:       "valid",
		Certificates: []string{base64.StdEncoding.EncodeToString(der)},
	}
	key := responses.CscKey{Status: "enabled", Algo: []string{cscECPublicKeyOID}}

	if parsed, perr := x509.ParseCertificate(der); perr == nil {
		cert.SubjectDN = parsed.Subject.String()
		cert.IssuerDN = parsed.Issuer.String()
		cert.SerialNumber = strings.ToUpper(parsed.SerialNumber.Text(16))
		cert.ValidFrom = parsed.NotBefore.UTC().Format("20060102150405Z")
		cert.ValidTo = parsed.NotAfter.UTC().Format("20060102150405Z")
		cert.QcStatements = cscQCStatements(parsed)
		if pub, ok := parsed.PublicKey.(*ecdsa.PublicKey); ok {
			key.Len = pub.Curve.Params().BitSize
		}
	}

	cscWriteJSON(w, http.StatusOK, responses.CscCredentialInfoResponse{
		Key:                key,
		Cert:               cert,
		SignatureQualifier: cscDefaultQualifier,
		Auth:               responses.CscAuth{Mode: "oauth2code"},
		SCAL:               "2",
		Multisign:          50,
	})
}

// CscSignHashHandler — POST /csc/v2/signatures/signHash. Signs each base64 digest
// with the ECC key and returns the base64 signatures in request order. Reads
// `hashes` (spec) or `hash` (alias). Encoding is DER by default;
// CSC_MOCK_SIG_ENCODING=P1363 emits raw r||s.
func CscSignHashHandler(privateKey *ecdsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isPostMethod(r) {
			cscError(w, http.StatusMethodNotAllowed, "invalid_request", "use POST")
			return
		}
		if privateKey == nil {
			cscError(w, http.StatusInternalServerError, "server_error", "signing key not loaded")
			return
		}

		var req requests.CscSignHashRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			cscError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
			return
		}

		digests := req.Digests()
		if len(digests) == 0 {
			cscError(w, http.StatusBadRequest, "invalid_request", "no hashes to sign")
			return
		}

		method := env.CscMockSigEncoding
		signatures := make([]string, len(digests))
		for i, h := range digests {
			sig, err := SignECCDigest(privateKey, h, method)
			if err != nil {
				cscError(w, http.StatusBadRequest, "invalid_request", fmt.Sprintf("hash %d: %s", i, err))
				return
			}
			signatures[i] = sig
		}

		log.Printf("csc-mock: signHash signed %d hash(es) using %s", len(digests), method)
		cscWriteJSON(w, http.StatusOK, responses.CscSignHashResponse{Signatures: signatures})
	}
}

// cscCertToDER returns the DER bytes of a certificate stored as PEM or as
// (possibly whitespace-wrapped) base64-DER.
func cscCertToDER(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if block, _ := pem.Decode([]byte(s)); block != nil {
		return block.Bytes, nil
	}
	return base64.StdEncoding.DecodeString(strings.Join(strings.Fields(s), ""))
}

// cscQCStatements extracts the QCStatement OIDs from a certificate's qcStatements
// extension, or nil when absent/unparseable.
func cscQCStatements(cert *x509.Certificate) []string {
	for _, ext := range cert.Extensions {
		if !ext.Id.Equal(oidQCStatements) {
			continue
		}

		var stmts []struct {
			ID   asn1.ObjectIdentifier
			Info asn1.RawValue `asn1:"optional"`
		}
		if _, err := asn1.Unmarshal(ext.Value, &stmts); err != nil {
			return nil
		}

		out := make([]string, 0, len(stmts))
		for _, s := range stmts {
			out = append(out, s.ID.String())
		}
		return out
	}
	return nil
}

// cscWriteJSON writes v as a JSON response with the given status.
func cscWriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("csc-mock: failed to encode response: %s", err)
	}
}

// cscError writes a CSC error body with the given HTTP status.
func cscError(w http.ResponseWriter, status int, code, desc string) {
	cscWriteJSON(w, status, responses.CscError{Error: code, ErrorDescription: desc})
}
