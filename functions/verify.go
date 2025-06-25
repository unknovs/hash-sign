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
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
)

func parseCertificate(certificateStr string) (*x509.Certificate, error) {
	certificateBytes, err := base64.StdEncoding.DecodeString(certificateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate: %v", err)
	}

	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return certificate, nil
}

func decodeBase64(s string) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}

	return bytes, nil
}

func verifyECDSASignature(pub *ecdsa.PublicKey, digestValue, signatureBytes []byte) error {
	var esig struct {
		R, S *big.Int
	}
	_, errAsn1 := asn1.Unmarshal(signatureBytes, &esig)
	if errAsn1 == nil {
		if !ecdsa.Verify(pub, digestValue, esig.R, esig.S) {
			return errors.New("ECDSA verification failed")
		}
	} else {
		keyBytes := (pub.Params().BitSize + 7) >> 3
		if len(signatureBytes) != 2*keyBytes {
			return errors.New("invalid ECDSA signature length")
		}
		r := new(big.Int).SetBytes(signatureBytes[:keyBytes])
		s := new(big.Int).SetBytes(signatureBytes[keyBytes:])
		if !ecdsa.Verify(pub, digestValue, r, s) {
			return errors.New("ECDSA verification failed")
		}
	}

	return nil
}

func VerifySignature(w http.ResponseWriter, r *http.Request) {
	var verifyBody requests.VerifyBody
	err := json.NewDecoder(r.Body).Decode(&verifyBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusUnprocessableEntity)
		return
	}

	certificate, err := parseCertificate(verifyBody.Certificate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signatureBytes, err := decodeBase64(verifyBody.SignatureValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid signature value: %v", err), http.StatusBadRequest)
		return
	}

	digestValue, err := decodeBase64(verifyBody.DigestValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid digest value: %v", err), http.StatusBadRequest)
		return
	}

	switch pub := certificate.PublicKey.(type) {
	case *rsa.PublicKey:
		err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, digestValue, signatureBytes)
	case *ecdsa.PublicKey:
		err = verifyECDSASignature(pub, digestValue, signatureBytes)
	default:
		http.Error(w, fmt.Sprintf("Unsupported public key type: %T", certificate.PublicKey), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify signature: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Signature is valid!")
}
