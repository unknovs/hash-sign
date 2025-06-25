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
