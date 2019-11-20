package main

import (
	"os"
	"log"
	"net/http"
	"strings"
	"io"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

var fs http.FileSystem

func main() {
	listeningIp := getEnv("HTTP_SERVER_LISTENING_IP", "0.0.0.0")
	listeningPort := getEnv("HTTP_SERVER_LISTENING_PORT", "80")
	fileDirectory := getEnv("HTTP_SERVER_DIRECTORY", "/static-data")
	prefix := getEnv("HTTP_SERVER_PREFIX", "/static/")

	fs = neuteredFileSystem{http.Dir(fileDirectory)}
	fileServer := http.FileServer(fs)
	if "" != prefix {
		fileServer = http.StripPrefix(prefix, fileServer)
	} else {
		prefix = "/"
	}

	fileServer = Handle404(fileServer, Fire404)

	http.Handle(prefix, fileServer)

	log.Println("Listening on " + listeningIp + ":" + listeningPort + "...")
	http.ListenAndServe(listeningIp + ":" + listeningPort, nil)
}

type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if s.IsDir() {
        index := strings.TrimSuffix(path, "/") + "/index.html"
        if _, err := nfs.fs.Open(index); err != nil {
            return nil, err
        }
    }

    return f, nil
}

type hijack404 struct {
    http.ResponseWriter
    R *http.Request
    Handle404 func (w http.ResponseWriter, r *http.Request) bool
}

func (h *hijack404) WriteHeader(code int) {
    if 404 == code && h.Handle404(h.ResponseWriter, h.R) {
        panic(h)
    }

    h.ResponseWriter.WriteHeader(code)
}

func Handle404(handler http.Handler, handle404 func (w http.ResponseWriter, r *http.Request) bool) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        hijack := &hijack404{ ResponseWriter:w, R: r, Handle404: handle404 }

        defer func() {
            if p:=recover(); p!=nil {
                if p==hijack {
                    return
                }
                panic(p)
            }
        }()

        handler.ServeHTTP(hijack, r)
    })
}

func Fire404(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusNotFound) // StatusNotFound = 404
	file, err := fs.Open("404.html")
	if err != nil {
	    w.Write([]byte("404 not found"))
	    return true;
	}
	_, err = io.Copy(w, file)
	if err != nil {
	    w.Write([]byte("404 not found"))
	}
	file.Close() // consider defer ^
	return true;
}