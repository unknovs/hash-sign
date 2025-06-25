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
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/routes/requests"
	"github.com/unknovs/hash-sign/routes/responses"
)

func HandleDigest(w http.ResponseWriter, r *http.Request) {

	if !isPostMethod(r) {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body
	var req requests.DigestSummaryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON request body: %v", err)
		return
	}

	// Get the digest from the request
	digest := req.DigestToCalculate

	if req.DigestToCalculate == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: Digest not provided")
		return
	}

	// Parse the query parameters
	query := r.URL.Query()
	hash := query.Get("hash")
	if hash == "" {
		hash = "sha256" // Default to sha256 if hash is not provided
	} else if hash != "sha256" && hash != "sha384" && hash != "sha512" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported hash algorithm: %s", hash)
		return
	}

	// Base64 to binary
	binaryDigest, err := base64.StdEncoding.DecodeString(digest)
	if err != nil {
		binaryDigest, err = base64.URLEncoding.DecodeString(digest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error decoding base64 digest: %v", err)
			return
		}
	}

	// Hash binary digest with selected algorithm
	var hashedDigest []byte
	var algorithm string
	switch hash {
	case "sha256":
		hashed := sha256.Sum256(binaryDigest)
		hashedDigest = hashed[:]
		algorithm = "sha256"
	case "sha384":
		hashed := sha512.Sum384(binaryDigest)
		hashedDigest = hashed[:]
		algorithm = "sha384"
	case "sha512":
		hashed := sha512.Sum512(binaryDigest)
		hashedDigest = hashed[:]
		algorithm = "sha512"
	}

	// Hex to base64
	digestSummary := base64.StdEncoding.EncodeToString(hashedDigest)
	URLSafeDigestSummary := base64.URLEncoding.EncodeToString(hashedDigest)

	response := responses.DigestSummary{
		DigestSummary:  digestSummary,
		UrLSafeSummary: URLSafeDigestSummary,
		Algorithm:      algorithm,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// spam a bit in log
	log.Printf("Digest summary of digest %v calculated with hash %s", digest, hash)

	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var writeErr error
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Printf("Error writing JSON response: %v", writeErr)
	}
}
