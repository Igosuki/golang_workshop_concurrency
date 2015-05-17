package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"path"
	// "path/filepath"
	"strings"
)

var chttp *http.ServeMux

func main() {
	log.SetLevel(log.DebugLevel)
	args := os.Args
	switch args[1] {
	case "pingpong":
		pingpong()
	case "rss":
		chttp = http.NewServeMux()
		chttp.Handle("/", http.FileServer(http.Dir("./public")))
		http.HandleFunc("/", HomeHandler)
		http.ListenAndServe(":8000", nil)
		// feed([]string{"blog.golang.org", "news.ycombinator.com/rss"})

	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") || r.URL.Path == "/" {
		log.Debug(path.Clean(r.URL.Path))
		chttp.ServeHTTP(w, r)
	} else {
		fmt.Fprintf(w, "HomeHandler")
	}
}
