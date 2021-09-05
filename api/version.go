package api

import "net/http"

func Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("1.0.0"))
}
