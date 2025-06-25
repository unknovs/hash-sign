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
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/unknovs/hash-sign/env"
)

func isPostMethod(r *http.Request) bool {
	return r.Method == http.MethodPost
}

// func isGetMethod(r *http.Request) bool {
// 	return r.Method == http.MethodGet
// }

func isGetMethod(r *http.Request) bool {
	return strings.TrimSpace(r.Method) == http.MethodGet
}

func APIKeyAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if env.ApiKey != "" && env.ApiKey != r.Header.Get("API-Key") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// func APIKeyAuthorization(apiKey string, next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if apiKey != "" && apiKey != r.Header.Get("API-Key") {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		next(w, r)
// 	}
// }

func CheckVolumeMounted(volumePath string) bool {
	// Attempt to obtain information about the volume without creating a directory
	if _, err := os.Stat(volumePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Volume is not mounted at %s", volumePath)
		} else {
			log.Printf("Error checking volume: %v", err)
		}
		return false
	}

	log.Printf("Volume is mounted at %s", volumePath)
	return true
}
