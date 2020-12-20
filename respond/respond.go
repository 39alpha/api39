package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"error"`
	Status  int    `json:"-"`
}

func NewErrorResponse(message string, status int) *ErrorResponse {
	return &ErrorResponse{message, status}
}

func MakeErrorResponse(err error, status int) *ErrorResponse {
	return NewErrorResponse(fmt.Sprintf("%q", err), status)
}

func (e ErrorResponse) WriteTo(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if blob, err := json.Marshal(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{ \"error\": \"failed to marshal an error response\" }"))
	} else {
		w.WriteHeader(e.Status)
		w.Write(blob)
	}
}

func WithErrorMessage(w http.ResponseWriter, message string, status int) {
	NewErrorResponse(message, status).WriteTo(w)
}

func WithError(w http.ResponseWriter, err error, status int) {
	MakeErrorResponse(err, status).WriteTo(w)
}

func ServerError(w http.ResponseWriter) {
	WithErrorMessage(w, "internal server error", http.StatusInternalServerError)
}

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if blob, err := json.Marshal(data); err != nil {
		WithError(w, err, http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(blob)
	}
}

func Raw(w http.ResponseWriter, data string) {
	RawWithCode(w, data, http.StatusOK)
}

func RawWithCode(w http.ResponseWriter, data string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}
