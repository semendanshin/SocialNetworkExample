package errorwrapper

import (
	"net/http"
)

func WriteWithError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + err + `"}`))
}

type HandlerWithError func(http.ResponseWriter, *http.Request) error

func WrapWithError(handlerFunc HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerFunc(w, r)
		if err != nil {
			WriteWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
}
