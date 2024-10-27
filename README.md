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

## Static frontend

Our reverse proxy handles anything not found and routes to the sveltekit
dev server, but what if we want to build the site and serve it statically
from our exe?

Right now I'll create a string to define how we want to serve the pages.
In the production build they will be embedded in the exe, but for now
I'll just declare a string in `main.go`:

```go
var frontend_mode = "static" // options 'static', 'dev', 'embedded' (default)
```

And for now (proxy code collapsed to `...`):

```go
	if frontend_mode == "proxy" {
        ...
	} else {
		fs := http.FileServer(http.Dir("./frontend/build"))
		http.Handle("/", fs)
	}
```

I guess svelte is a little weird as it's meant to support various platforms
and server-side rendering.   To build a static site we have to follow
[this page](https://svelte.dev/docs/kit/adapter-static)

First, install `@sveltejs/adapter-static` as a development dependency
in our `frontend` directory with the svelte app:

    npm i -D @sveltejs/adapter-static

Then use that instead of the auto one in `frontend/svelte.config.js`:

```js
// import adapter from '@sveltejs/adapter-auto';
import adapter from '@sveltejs/adapter-static';
```

And we have to tell `frontend/src/routes/+layout.js` to prerender everything:

```js
export const prerender = true;
export const ssr = false;
```

Now I can run `npm run build` in the `frontend` directory and get the output
in `frontend/build`.   Tested with python `python -m http.server` in that
directory and serving.

And running the go app with `go run .` works too, with one caveat.
The file system server serves things fine and serves `index.html`
if no route is specified, but I added other routes to the svelte app.
These work fine once the app is up because it does everything on the
client, but if you go directly to `http://localhost:7000/about`
you will get a 404 error.

This commit is big enough for now because I added code for routes to
the svelte app, so I'm committing now.  I learned a bit about the routing
in svelte along the way.   Had to add a `frontend/src/routes/+layout.svelte`
file with a `<slot></slot>` element to get content to appear
on every page.   

That's what the `+layout.js` file means too I think, that runs for every
page and returns that it should be prerendered with ssr disabled.  I guess
this could go in any specific route as well, kinda cool.

## Embedding frontend

We can use the `embed` package to embed the spa files into our app.  I added
the file `embedHandler.go` to contain the code.   Note the `//go.embed`
comment actually works a directive that tells the compiler what files to
embed in the executable.   The other functions are to convert paths, otherwise
we would need to go to 'http://localhost:7000/frontend/build' to see the
root page of our spa.  I also check for an error opening an embedded file
(which happens if it isn't found) and return the root `index.html` instead
so we can go directly to sub-routes from typing in the address.

```go
package main

import (
	"embed"
	"io/fs"
	"path"
)

//go:embed frontend/build/*
var frontendEmbedded embed.FS

type subdirFS struct {
	embed.FS
	subdir string
}

func (s subdirFS) Open(name string) (fs.File, error) {
	file, err := s.FS.Open(path.Join(s.subdir, name))
	if err == nil {
		return file, nil
	}
	file, err = s.FS.Open(path.Join(s.subdir, "index.html"))
	return file, err
}

var frontendFs = subdirFS{frontendEmbedded, "frontend/build"}
```

And serving in `main.go`:

```go
http.Handle("/", http.FileServer(http.FS(frontendFs)))
```

This shows the power and flexibility of GO's interface system.
`http.FileServer` requires an FS object, those have a lot
of methods.   The `http.FS()` method takes an object
that implements the [http.FileSystem](https://pkg.go.dev/net/http#FileSystem)
interface, which is only `Open(name string) File, error`.
`embed.FS` lets us access embedded files easily and we can check
for errors and return a different file.

We should probably do something similar with the static server.
