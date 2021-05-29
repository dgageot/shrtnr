package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/karrick/godirwalk"
)

//go:embed static/results.html.tmpl
var nearestResultTmpl string

func main() {
	r := mux.NewRouter()
	r.HandleFunc(`/`, index).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc(`/{key:[a-zA-Z][a-zA-Z0-9\-\_]*}`, link).Methods(http.MethodGet, http.MethodHead)
	log.Fatal(http.ListenAndServe(":8888", r))
}

func index(w http.ResponseWriter, r *http.Request) {
	key := r.Host
	doRedirect(key, w, r)
}

func link(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	doRedirect(key, w, r)
}

func doRedirect(key string, w http.ResponseWriter, r *http.Request) {
	url, err := urlForKey(key)
	switch {
	case os.IsNotExist(err):
		printMatches(key, w)

	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)

	default:
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func urlForKey(k string) (string, error) {
	path := filepath.Join("links", k)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	url := string(buf)
	url = strings.TrimSpace(string(buf))
	return url, nil
}

func printMatches(k string, w http.ResponseWriter) {
	matches, err := findMatches(k)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text.html")

	t := template.Must(template.New("tmpl").Parse(nearestResultTmpl))
	t.Execute(w, matches)
}

func findMatches(notFound string) (map[string]string, error) {
	matches := map[string]string{}

	scanner, err := godirwalk.NewScanner("links")
	if err != nil {
		return nil, fmt.Errorf("cannot scan directory: %w", err)
	}

	for scanner.Scan() {
		dirent, err := scanner.Dirent()
		if err != nil {
			return nil, fmt.Errorf("cannot get dirent: %w", err)
		}

		if dirent.IsRegular() {
			key := dirent.Name()
			if levenshteinDistance(key, notFound) <= 2 {
				url, err := urlForKey(key)
				if err != nil {
					return nil, err
				}

				matches[key] = url
			}
		}
	}

	return matches, nil
}
