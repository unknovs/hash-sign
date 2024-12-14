package functions

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unknovs/hash-sign/routes/requests"
	"github.com/unknovs/hash-sign/routes/responses"
)

func generateJWT(req requests.JWTRequest) (string, error) {

	// Read private key from environment variable

	privateKeyPEM := os.Getenv("JWT_SIGNING_KEY")

	// Read the entire file content
	keyBytes, err := os.ReadFile(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %v", err)
	}

	// Decode PEM encoded private key
	if privateKeyPEM == "" {
		return "", fmt.Errorf("JWT_SIGNING_KEY environment variable is not set")
	}

	block, rest := pem.Decode(keyBytes)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block. Decoded content: %v, Remaining content: %v", block, string(rest))
	}

	// Parse the private key (now using ParsePKCS8PrivateKey for PKCS8 format)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Validate input
	if req.Issuer == "" || req.Audience == "" || req.Subject == "" {
		return "", fmt.Errorf("issuer, audience, and subject must be provided")
	}

	// Generate a random JTI (JWT ID)
	jtiBytes := make([]byte, 18) // 24 base64 chars = 18 raw bytes
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", fmt.Errorf("failed to generate JTI: %v", err)
	}
	jti := base64.StdEncoding.EncodeToString(jtiBytes)

	// Create claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": req.Issuer,                      // Issuer
		"aud": req.Audience,                    // Audience
		"sub": req.Subject,                     // Subject
		"iat": now.Unix(),                      // Issued At
		"exp": now.Add(5 * time.Minute).Unix(), // Expiration Time (5 minutes)
		"jti": jti,                             // JWT ID
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func JwtGenerateHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it's a POST request
	if !isPostMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	// Parse JSON request body
	var jwtRequest requests.JWTRequest
	err := json.NewDecoder(r.Body).Decode(&jwtRequest)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Empty request body", http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		}
		return
	}

	// Generate JWT
	tokenString, err := generateJWT(jwtRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("JWT generation successful")

	// Create response
	response := responses.JWTResponse{Token: tokenString}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send JSON response
	json.NewEncoder(w).Encode(response)
}
