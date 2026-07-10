package main

import (
	"embed"
	"net/http"
)

//go:embed all:static/*
var gui embed.FS

func FileServer(w http.ResponseWriter, r *http.Request) {
	path := "static" + r.URL.Path
	if path[len(path)-1] == '/' {
		path += "index.html"
	}
	http.ServeFileFS(w, r, gui, path)
}
