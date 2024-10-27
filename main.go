package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/JasonGoemaat/go-spa/api"
)

func main() {
	fmt.Println("Hello, World!")

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, World!"))
	// })

	remote, _ := url.Parse("http://localhost:5173")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = remote.Host
		proxy.ServeHTTP(w, r)
	})

	// add our api call, anything not handled elsewhere is handled
	// by the "/" above
	http.HandleFunc("/api/hello", api.Hello)

	fmt.Println("Listening on localhost port 7000, proxying to localhost:5173")
	err := http.ListenAndServe("localhost:7000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
