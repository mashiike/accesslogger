package accesslogger

import (
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	FirstWriteTime time.Time
	LastWriteTime  time.Time
	StatusCode     int
	BodyByteSent   int
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if w.FirstWriteTime.IsZero() {
		w.FirstWriteTime = Clock()
	}
	w.LastWriteTime = Clock()
	w.BodyByteSent += n
	return n, err
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
