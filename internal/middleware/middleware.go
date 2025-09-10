package middleware

import (
	"net/http"
	"strings"
)

func CheckApplicationJSONContentType(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		h(w, r)
	}
}

func CheckPlainTextContentType(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" && !strings.HasPrefix(contentType, "text/plain") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		h(w, r)
	}
}
