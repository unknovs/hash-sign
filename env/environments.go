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

package env

import (
	"log"
	"os"
	"strings"
)

var (
	PemFile          = os.Getenv("PEM_FILE")
	EcPemFile        = os.Getenv("EC_PEM_FILE")
	ApiKey           = os.Getenv("API_KEY")
	RsaAuthCert      = os.Getenv("RSA_AUTH_CERT")
	RsaSigningCert   = os.Getenv("RSA_SIGN_CERT")
	EcdsaAuthCert    = os.Getenv("ECDSA_AUTH_CERT")
	EcdsaSigningCert = os.Getenv("ECDSA_SIGN_CERT")
	jwtSigningKey    = os.Getenv("JWT_SIGNING_KEY")
)

// getEnvOrSecret reads the environment variable or Docker secret file content.
func getEnvOrSecret(varName string) string {
	value := os.Getenv(varName)
	if strings.HasPrefix(value, "/run/secrets/") {
		data, err := os.ReadFile(value)
		if err != nil {
			log.Printf("unable to read secret file %s: %v", value, err)
			return ""
		}
		return strings.TrimSpace(string(data))
	}
	return value
}
