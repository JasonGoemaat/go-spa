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
