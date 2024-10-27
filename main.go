package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/JasonGoemaat/go-spa/api"
)

var frontend_mode = "embed" // options 'static', 'dev', 'embed' (default)

func main() {
	fmt.Println("Hello, World!")

	if frontend_mode == "dev" {
		fmt.Println("Using 'dev' mode!")
		fmt.Println("Proxying requests to http://localhost:5173")
		remote, _ := url.Parse("http://localhost:5173")
		proxy := httputil.NewSingleHostReverseProxy(remote)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			r.Host = remote.Host
			proxy.ServeHTTP(w, r)
		})
	} else if frontend_mode == "static" {
		fmt.Println("Using 'static' mode!")
		fmt.Println("Serving files from ./frontend/build")
		http.Handle("/", http.FileServer(http.Dir("./frontend/build")))
	} else if frontend_mode == "embed" {
		fmt.Println("Using 'embed' mode!")
		fmt.Println("Using embedded files from ./frontend/build")
		http.Handle("/", http.FileServer(http.FS(frontendFs)))
	} else {
		panic(fmt.Sprintf("Unknown frontend_mode: '%s'", frontend_mode))
	}

	// add our api call, anything not handled elsewhere is handled
	// by the "/" above
	http.HandleFunc("/api/hello", api.Hello)

	fmt.Println("Listening on localhost port 7000")
	err := http.ListenAndServe("localhost:7000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
