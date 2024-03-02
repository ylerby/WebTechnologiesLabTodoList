package response

import "net/http"

const (
	headerKey   = "Content-Type"
	contentType = "application/json"
)

// CorrectResponseWriter todo: logger or return error interface
func CorrectResponseWriter(w http.ResponseWriter, data []byte, statusCode int) error {
	w.Header().Set(headerKey, contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func ErrorResponseWriter(w http.ResponseWriter, data []byte, statusCode int) error {
	w.Header().Set(headerKey, contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		return err
	}

	return nil
}
