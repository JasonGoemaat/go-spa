package api

import (
	"fmt"
	"net/http"
)

// /iisintegration URL is called when IIS shuts down the app pool
func IISshutdown(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"message\":\"Hello, world!\"}")
}
