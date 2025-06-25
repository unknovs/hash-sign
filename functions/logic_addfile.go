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

func decodeRequest(r *http.Request) (requests.Request, error) {
	var req requests.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Can't decode request, getting error: %v", err)
	}
	return req, err
}

func getEmptyAsiceReader(req requests.Request) (*zip.Reader, error) {
	SignedEmptyAsiceBytes, err := base64.StdEncoding.DecodeString(req.EmptyAsice)
	if err != nil {
		log.Printf("Can't decode Base64-encoded empty ASiC-E from request, getting error: %v", err)
		return nil, err
	}

	emptyAsiceReader, err := zip.NewReader(bytes.NewReader(SignedEmptyAsiceBytes), int64(len(SignedEmptyAsiceBytes)))
	if err != nil {
		log.Printf("Error reading decoded ASic-E: %s", err.Error())
	}
	return emptyAsiceReader, err
}

func createNewAsiceFile() (*os.File, *zip.Writer, error) {
	newAsiceFile, err := os.CreateTemp("", "newZip")
	if err != nil {
		log.Printf("Error creating temporary file for new ASiC-E file: %v", err)
		return nil, nil, err
	}

	newAsiceWriter := zip.NewWriter(newAsiceFile)
	return newAsiceFile, newAsiceWriter, nil
}

func addFilesToArchive(req requests.Request, emptyAsiceReader *zip.Reader, newAsiceWriter *zip.Writer) error {
	for _, file := range req.SignedFiles {
		signedFile := requests.SignedFile(file) // Convert File to SignedFile
		if err := addFileToArchive(newAsiceWriter, signedFile); err != nil {
			return err
		}
	}

	for _, file := range emptyAsiceReader.File {
		if file.Mode().IsDir() {
			continue
		}

		if err := addFileToArchiveFromReader(newAsiceWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func writeResponse(w http.ResponseWriter, r *http.Request, newAsiceFile *os.File) error {
	newAsiceFileBytes, err := os.ReadFile(newAsiceFile.Name())
	if err != nil {
		log.Printf("Error reading newAsiceFile: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	fileType := r.URL.Query().Get("type")
	log.Println("Provided files added to ASiC-E container")

	if fileType == "binary" {
		w.Header().Set("Content-Type", "application/zip")
		_, err := w.Write(newAsiceFileBytes)
		if err != nil {
			log.Printf("Error writing newAsiceFileBytes: %v", err)
		}
		return err
	} else if fileType == "base64" {
		response := base64.StdEncoding.EncodeToString(newAsiceFileBytes)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"packedAsice": response,
		})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return err
	} else {
		_, err := w.Write(newAsiceFileBytes)
		if err != nil {
			log.Printf("Error writing newAsiceFileBytes: %v", err)
		}
		return err
	}
}

func addFileToArchive(archive *zip.Writer, file requests.SignedFile) error {
	decodedFileBytes, err := base64.StdEncoding.DecodeString(file.EncodedFile)
	if err != nil {
		log.Printf("Error decoding file %s: %v", file.FileName, err)
		return err
	}

	newFileWriter, err := archive.Create(file.FileName)
	if err != nil {
		log.Printf("Error creating new file %s in the ASiC-E archive: %v", file.FileName, err)
		return err
	}

	_, err = newFileWriter.Write(decodedFileBytes)
	if err != nil {
		log.Printf("Error writing file %s to the ASiC-E archive: %v", file.FileName, err)
		return err
	}

	return nil
}

func addFileToArchiveFromReader(archive *zip.Writer, file *zip.File) error {
	if file.Mode().IsDir() {
		return nil
	}

	newFileWriter, err := archive.Create(file.Name)
	if err != nil {
		log.Printf("Error creating new file in ASiC-E archive: %v", err)
		return err
	}

	fileReader, err := file.Open()
	if err != nil {
		log.Printf("Error opening file in ASiC-E archive: %v", err)
		return err
	}
	defer fileReader.Close()

	_, err = io.Copy(newFileWriter, fileReader)
	if err != nil {
		log.Printf("Error copying file to new ASiC-E archive: %v", err)
		return err
	}

	return nil
}
