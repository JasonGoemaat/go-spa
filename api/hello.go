package api

import (
	"fmt"
	"net/http"

	"github.com/JasonGoemaat/go-spa/mylog"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	mylog.Log("API:Hello()")
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"message\":\"Hello, world!\"}")
}
