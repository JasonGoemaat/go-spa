# go-spa

## Initial Commit 

Created repo on github with MIT license and go `.gitignore` specified.

## Hello, World!

Initialized a go mod with this:

    go mod init github.com/JasonGoemaat/go-spa

Wrote `main.go` file to print 'Hello, World!' to console.

## Simple Server

Listen on localhost port 7000 and serve 'Hello, World!'
to requests.

## Svelte app

Created a svelte app in the `frontend` directory.

    npx sv create frontend

Selected `SvelteKit minimal`

Selected `Yes, using Typescript syntax`

Did not select additions

Selected `npm`

This seems to install dependencies automatically.
So I change to the directory and check it out:

    cd frontend
    npm run dev

And content is served on `http://localhost:5173/`

## Reverse proxy

Goal is to add a handler to reverse proxy the svelte app
to our go app running on port 7000.

Well, that was pretty easy...  Go includes a reverse proxy
in it's standard library using a handler func like
everything else.

```go
remote, _ := url.Parse("http://localhost:5173")
proxy := httputil.NewSingleHostReverseProxy(remote)
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    r.Host = remote.Host
    proxy.ServeHTTP(w, r)
})
```

## Hello API

I added an 'api' directory which will be a package and a
'hello.go' file in it with a function to return JSON:

```go
package api

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"message\":\"Hello, world!\"}")
}
```

And we can add a handler for the `/api/hello` route and the existing handler
serving "/" will handle anything else.

```go
http.HandleFunc("/api/hello", api.Hello)
```
