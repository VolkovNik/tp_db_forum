package utils

import (
	"TPForum/internal/pkg/domain"
	"encoding/json"
	"net/http"
)

func WriteResponseError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(domain.NewError(msg))
	if err != nil {
		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
	}

	if _, err = w.Write(data); err != nil {
		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
	}

}

func WriteResponse(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
	}
}
