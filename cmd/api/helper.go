package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonRequest struct {
	Error   bool   `josn:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1048576
	reader := http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("should only have one json block")
	}
	return nil
}

func writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	jData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	// this needs to be removed after testing, it should not allow access from webbrowser
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.WriteHeader(status)

	_, err = w.Write(jData)

	if err != nil {
		return err
	}

	return nil
}

func errorJson(w http.ResponseWriter, message string, status ...int) error {
	errorStatus := http.StatusBadRequest
	if len(status) > 0 {
		errorStatus = status[0]
	}

	payLoad := jsonRequest{
		Error:   true,
		Message: message,
	}

	return writeJson(w, errorStatus, payLoad)
}
