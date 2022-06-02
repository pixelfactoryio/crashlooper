package handlers

import (
	"net/http"
)

type defaultHandler struct{}

// NewDefaultHandler returns a new defaultHandler instance.
func NewDefaultHandler() http.Handler {
	return &defaultHandler{}
}

// ServeHTTP respond with default index.
func (h *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(`
	<html>
		<head>
			<title>CrashLooper</title>
		</head>
		<body>
			<h1>CrashLooper</h1>
			<p><a href='/shutdown'>Shutdown</a></p>
		</body>
	</html>`,
	))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
