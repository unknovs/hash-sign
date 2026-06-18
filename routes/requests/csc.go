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

package requests

// CscSignHashRequest is the body of POST /csc/v2/signatures/signHash. The CSC v2
// field is `hashes`; `hash` is accepted as a fallback alias. `SAD` is optional
// (the oauth2code credential token is the authorization).
type CscSignHashRequest struct {
	CredentialID  string   `json:"credentialID"`
	SAD           string   `json:"SAD,omitempty"`
	Hashes        []string `json:"hashes"`
	Hash          []string `json:"hash,omitempty"`
	SignAlgo      string   `json:"signAlgo"`
	OperationMode string   `json:"operationMode,omitempty"`
}

// Digests returns the digests to sign, preferring `hashes` and falling back to
// the `hash` alias.
func (r CscSignHashRequest) Digests() []string {
	if len(r.Hashes) > 0 {
		return r.Hashes
	}
	return r.Hash
}
