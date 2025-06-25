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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	"github.com/unknovs/hash-sign/env"
	"github.com/unknovs/hash-sign/routes/responses"
)

func certificatesExist() bool {
	return env.RsaAuthCert != "" || env.RsaSigningCert != "" || env.EcdsaAuthCert != "" || env.EcdsaSigningCert != ""
}

func CheckCertificates() error {
	certs := map[string]string{
		"RsaAuthCert":      env.RsaAuthCert,
		"RsaSigningCert":   env.RsaSigningCert,
		"EcdsaAuthCert":    env.EcdsaAuthCert,
		"EcdsaSigningCert": env.EcdsaSigningCert,
	}

	for name, cert := range certs {
		if cert == "" {
			fmt.Printf("%s is not provided\n", name)
			continue
		}

		pemCert := fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", cert)
		block, _ := pem.Decode([]byte(pemCert))
		if block == nil {
			return fmt.Errorf("failed to parse certificate PEM for %s", name)
		}

		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate for %s: %v", name, err)
		}

		switch name {
		case "RsaAuthCert", "RsaSigningCert":
			if _, ok := certificate.PublicKey.(*rsa.PublicKey); !ok {
				return fmt.Errorf("provided %s is not an RSA public key", name)
			}

		case "EcdsaAuthCert", "EcdsaSigningCert":
			if _, ok := certificate.PublicKey.(*ecdsa.PublicKey); !ok {
				return fmt.Errorf("provided %s is not an ECDSA public key", name)
			}

		}
	}

	return nil
}

func getCertificatesResponse(key, certType string, r *http.Request) (responses.CertificatesResponse, error) {
	response := responses.CertificatesResponse{}

	switch {
	case key == "":
		return handleBothCertificates(certType, r)
	case key == "rsa":
		return handleRSACertificates(certType, r)
	case key == "ecdsa":
		return handleECDSACertificates(certType, r)
	default:
		return response, fmt.Errorf("invalid 'key' parameter")
	}
}

func handleRSACertificates(certType string, r *http.Request) (responses.CertificatesResponse, error) {
	response := responses.CertificatesResponse{}

	switch certType {
	case "":
		if env.RsaAuthCert == "" && env.RsaSigningCert == "" {
			return response, fmt.Errorf("RSA certificates not found")
		}
		response.RsaAuthenticationCertificate = env.RsaAuthCert
		response.RsaSigningCertificate = env.RsaSigningCert
		log.Printf("%s?key=rsa responded", r.URL.Path)
	case "auth":
		if env.RsaAuthCert == "" {
			return response, fmt.Errorf("RSA Authentication Cert not found")
		}
		response.RsaAuthenticationCertificate = env.RsaAuthCert
		log.Printf(logMessage, r.URL.Path, "rsa", certType)
	case "sign":
		if env.RsaSigningCert == "" {
			return response, fmt.Errorf("RSA Signing Cert not found")
		}
		response.RsaSigningCertificate = env.RsaSigningCert
		log.Printf(logMessage, r.URL.Path, "rsa", certType)
	default:
		return response, fmt.Errorf("invalid 'type' parameter for key 'rsa'")
	}

	return response, nil
}

func handleECDSACertificates(certType string, r *http.Request) (responses.CertificatesResponse, error) {
	response := responses.CertificatesResponse{}

	switch certType {
	case "":
		if env.EcdsaAuthCert == "" && env.EcdsaSigningCert == "" {
			return response, fmt.Errorf("ECDSA certificates not found")
		}
		response.EcdsaAuthenticationCertificate = env.EcdsaAuthCert
		response.EcdsaSigningCertificate = env.EcdsaSigningCert
		log.Printf("%s?key=ecdsa responded", r.URL.Path)
	case "auth":
		if env.EcdsaAuthCert == "" {
			return response, fmt.Errorf("ECDSA Authentication Cert not found")
		}
		response.EcdsaAuthenticationCertificate = env.EcdsaAuthCert
		log.Printf(logMessage, r.URL.Path, "ecdsa", certType)
	case "sign":
		if env.EcdsaSigningCert == "" {
			return response, fmt.Errorf("ECDSA Signing Cert not found")
		}
		response.EcdsaSigningCertificate = env.EcdsaSigningCert
		log.Printf(logMessage, r.URL.Path, "ecdsa", certType)
	default:
		return response, fmt.Errorf("invalid 'type' parameter for key 'ecdsa'")
	}

	return response, nil
}

func handleBothCertificates(certType string, r *http.Request) (responses.CertificatesResponse, error) {
	response := responses.CertificatesResponse{}
	switch certType {
	case "":
		response.RsaAuthenticationCertificate = env.RsaAuthCert
		response.RsaSigningCertificate = env.RsaSigningCert
		response.EcdsaAuthenticationCertificate = env.EcdsaAuthCert
		response.EcdsaSigningCertificate = env.EcdsaSigningCert
		log.Printf("%s responded", r.URL.Path)
	default:
		return response, fmt.Errorf("cannot use 'type' without 'key' in querry")
	}

	return response, nil
}
