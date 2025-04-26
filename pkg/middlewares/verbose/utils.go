package verbose

import "net/http"

type responseData struct {
	status int
	size   int
}

type verboseResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *verboseResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *verboseResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
