package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/JasonGoemaat/go-spa/mylog"
)

//go:embed frontend/build/*
var frontendEmbedded embed.FS

type subdirFS struct {
	embed.FS
	subdir string
}

var main_counter = 0

func (s subdirFS) Open(name string) (fs.File, error) {
	main_counter++
	counter := main_counter
	str := fmt.Sprintf("subdirFS(%d) - looking for '%s'", counter, name)
	mylog.Log(str)

	// should use http.StripPrefix in creating handler maybe?
	if strings.HasPrefix(name, "go-spa") {
		name = name[6:]
		str = fmt.Sprintf("subdirFS(%d) - removed go-spa: '%s'", counter, name)
		mylog.Log(str)
	}

	path1 := path.Join(s.subdir, name)
	str = fmt.Sprintf("subdirFS(%d) - path1 is '%s'", counter, path1)
	mylog.Log(str)

	file, err := s.FS.Open(path1)
	if err == nil {
		str = fmt.Sprintf("subdirFS(%d) - found, returning...", counter)
		mylog.Log(str)
		return file, nil
	}
	file, err = s.FS.Open(path.Join(s.subdir, "index.html"))
	if err != nil {
		str = fmt.Sprintf("subdirFS(%d) - error with index.html: %s", counter, err.Error())
		mylog.Log(str)
	} else {
		str = fmt.Sprintf("subdirFS(%d) - returning index.html", counter)
		mylog.Log(str)
	}
	return file, err
}

var frontendFs = subdirFS{frontendEmbedded, "frontend/build"}
