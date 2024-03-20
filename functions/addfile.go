package functions

import (
	"log"
	"net/http"
	"os"
)

func HandleAddFileToAsiceRequest(w http.ResponseWriter, r *http.Request) {
	if !isPostMethod(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !CheckVolumeMounted("/tmp") {
		log.Println("Volume is not available or mounted. asice/addFile method is not available")
		return
	}

	req, err := decodeRequest(r)
	if err != nil {
		http.Error(w, "Can't decode request, getting error", http.StatusBadRequest)
		return
	}

	emptyAsiceReader, err := getEmptyAsiceReader(req)
	if err != nil {
		http.Error(w, "Error reading decoded ASic-E:", http.StatusBadRequest)
		return
	}

	newAsiceFile, newAsiceWriter, err := createNewAsiceFile()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(newAsiceFile.Name())

	if err := addFilesToArchive(req, emptyAsiceReader, newAsiceWriter); err != nil {
		return
	}

	if err := newAsiceWriter.Close(); err != nil {
		log.Printf("Error closing newAsiceWriter: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := writeResponse(w, r, newAsiceFile); err != nil {
		return
	}
}
