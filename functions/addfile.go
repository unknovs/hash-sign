package functions

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unknovs/hash-sign/routes/requests"
)

func HandleAddFileToAsiceRequest(w http.ResponseWriter, r *http.Request) {

	if !isPostMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check if the volume is mounted
	volumePath := "/tmp" 
	if !CheckVolumeMounted(volumePath) {
		log.Println("Volume is not available or mounted. asice/addFile method is not available")
	}
	
	// Lest take a Request struct to hold the decoded JSON data
	var req requests.Request

	// Decode the request JSON data and store it in the Request struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Log an error message with the details of the error
		log.Printf("Can't decode request, getting error: %v", err)
		// Lets spam some error details to the response
		http.Error(w, "Can't decode request, getting error", http.StatusBadRequest)
		return
	}

	// Decode the Base64-encoded signed empty (without signed files) ASiC-E ZIP bytes from the request
	SignedEmptyAsiceBytes, err := base64.StdEncoding.DecodeString(req.EmptyAsice)

	// Check if an error occurred while decoding the bytes
	if err != nil {
		// Log an error message with the details of the error
		log.Printf("Can't decode Base64-encoded empty ASiC-E from request, getting error: %v", err)
		// Lets spam some error details to the response
		http.Error(w, "Can't decode Base64-encoded empty ASiC-E from request, getting error", http.StatusBadRequest)
		return
	}

	// Lets try to read our decoded ASic-E ZIP hiding in SignedEmptyAsiceBytes byte slice.
	emptyAsiceReader, err := zip.NewReader(bytes.NewReader(SignedEmptyAsiceBytes), int64(len(SignedEmptyAsiceBytes)))

	if err != nil {
		log.Printf("Error reading decoded ASic-E: %s", err.Error())
		http.Error(w, "Error reading decoded ASic-E:", http.StatusBadRequest)
		return
	}

	newAsiceFile, err := os.CreateTemp("", "newZip")
	if err != nil {
		log.Printf("Error creating temporary file for new ASiC-E file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(newAsiceFile.Name())

	newAsiceWriter := zip.NewWriter(newAsiceFile)

	for _, file := range req.SignedFiles {
		decodedFileBytes, err := base64.StdEncoding.DecodeString(file.EncodedFile)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newFileWriter, err := newAsiceWriter.Create(file.FileName)
		if err != nil {
			log.Printf("Error creating new file %s in the ASiC-E archive: %v", file.FileName, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = newFileWriter.Write(decodedFileBytes)
		if err != nil {
			log.Printf("Error writing file %s to the ASiC-E archive: %v", file.FileName, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	for _, file := range emptyAsiceReader.File {
		if file.Mode().IsDir() {
			continue
		}

		newFileWriter, err := newAsiceWriter.Create(file.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error creating new file in ASiC-E archive: %v", err)
			return
		}

		fileReader, err := file.Open()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error opening file in ASiC-E archive: %v", err)
			return
		}

		_, err = io.Copy(newFileWriter, fileReader)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error copying file to new ASiC-E archive: %v", err)
			return
		}
	}

	err = newAsiceWriter.Close()
	if err != nil {
		log.Printf("Error closing newAsiceWriter: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newAsiceFileBytes, err := os.ReadFile(newAsiceFile.Name())
	if err != nil {
		log.Printf("Error reading newAsiceFile: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response interface{}
	fileType := r.URL.Query().Get("type")

	log.Println("Provided files added to ASiC-E container")

	if fileType == "binary" {
		w.Header().Set("Content-Type", "application/zip")
		_, err := w.Write(newAsiceFileBytes)
		if err != nil {
			log.Printf("Error writing newAsiceFileBytes: %v", err)
		}

		return
	} else if fileType != "" {
		// Only allow "binary" or "base64"
		if fileType != "base64" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response = base64.StdEncoding.EncodeToString(newAsiceFileBytes)
		w.Header().Set("Content-Type", "application/json")
	} else {
		// Return the binary data with no headers set
		_, err := w.Write(newAsiceFileBytes)
		if err != nil {
			log.Printf("Error writing newAsiceFileBytes: %v", err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"packedAsice": response,
	})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
