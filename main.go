package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	fmt.Println("Listening on localhost port 7000, proxying to localhost:5173")
	http.ListenAndServe("localhost:7000", nil)
}
