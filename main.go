package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/JasonGoemaat/go-spa/api"
	"github.com/JasonGoemaat/go-spa/mylog"
)

var frontend_mode = "dev" // options 'static', 'dev', 'embed' (default)

var URL_PREFIX = ""

func main() {
	mylog.Log("--------------------------------------------------------------------------------")
	mylog.Log("main()")
	mylog.Log("Hello, World!")

	mylog.ShowEnvironmentVariables()
	// we can use `ASPNETCORE_PORT` for the port we listen on, picked by IIS
	// we can use `ASPNETCORE_APPL_PATH` env var for the root of our application ('/GO-SPA' when I have it as 'go-spa' under my Default Web Site)
	// NOTE: check `frontend/svelte.config.js`` where we set the path, this needs to be set when building (maybe use .env file?)
	// NOTE: check `frontend/src/routes/+layout.svelte`
	var APPL_PATH = os.Getenv("ASPNETCORE_APPL_PATH")
	URL_PREFIX = strings.ToLower(APPL_PATH) // this is upper-case for iis, but we compare with lower-case, also is conveniently lacking final '/'
	mylog.Log(fmt.Sprintf("URL_PREFIX: '%s'", URL_PREFIX))

	// add our api call, anything not handled elsewhere is handled
	// by the "/" above
	http.HandleFunc(URL_PREFIX+"/api/hello", api.Hello)
	http.HandleFunc(URL_PREFIX+"/iisintegration", api.IISshutdown)

	if frontend_mode == "dev" {
		mylog.Log("Using 'dev' mode!")
		mylog.Log("Proxying requests to http://localhost:5173/")
		remote, _ := url.Parse("http://localhost:5173/")
		proxy := httputil.NewSingleHostReverseProxy(remote)

		// NOTE: http.StripPrefix requires a Handler, not a handler func?
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			mylog.Log(fmt.Sprintf("Proxying %s", r.URL.Path))
			r.Host = remote.Host
			proxy.ServeHTTP(w, r)
		})
	} else if frontend_mode == "static" {
		mylog.Log("Using 'static' mode!")
		mylog.Log("Serving files from ./frontend/build")
		// http.Handle("/", http.FileServer(http.Dir("./frontend/build")))
		http.Handle("/", http.StripPrefix(URL_PREFIX, http.FileServer(http.FS(frontendFs))))
	} else if frontend_mode == "embed" {
		mylog.Log("Using 'embed' mode!")
		mylog.Log("Using embedded files from ./frontend/build")
		// http.Handle("/", http.FileServer(http.FS(frontendFs)))
		http.Handle("/", http.StripPrefix(URL_PREFIX, http.FileServer(http.FS(frontendFs))))
	} else {
		panic(fmt.Sprintf("Unknown frontend_mode: '%s'", frontend_mode))
	}

	port := "8080"
	if os.Getenv("ASPNETCORE_PORT") != "" { // get enviroment variable that set by ACNM
		port = os.Getenv("ASPNETCORE_PORT")
	}

	address := fmt.Sprintf(":%s", port)
	fmt.Println("Listening on", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err)
	}
}
