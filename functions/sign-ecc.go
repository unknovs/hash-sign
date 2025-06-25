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
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
)

func SigningHandlerEC(privateKey *ecdsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateRequest(w, r, privateKey) {
			return
		}

		var signEcdsa requests.SignEcdsa
		if !decodeJSON(w, r, &signEcdsa) {
			return
		}

		hashBytes, err := decodeHash(signEcdsa)
		if err != nil {
			http.Error(w, "Failed to decode hash from base64", http.StatusBadRequest)
			return
		}

		signatureR, signatureS, err := signHash(privateKey, hashBytes)
		if err != nil {
			http.Error(w, "Error signing hash", http.StatusInternalServerError)
			return
		}

		signatureMethod := getSignatureMethod(r)
		signature, err := encodeSignature(signatureMethod, signatureR, signatureS, privateKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sendResponse(w, signEcdsa, signature, signatureMethod)
	}
}
