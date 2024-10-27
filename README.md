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
