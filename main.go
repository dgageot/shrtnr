package main

import (
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

const nearestResultTmpl = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Nearest matches</title>
	</head>
	<body>
		<h1>Nearest matches:</h1>
		{{range $k,$v := .}}
		<p><a href="{{$v}}">go/{{$k}}</a></p>
		{{end}}
	</body>
</html>`

func main() {
	r := mux.NewRouter()
	r.HandleFunc(`/`, index).Methods("HEAD", "GET")
	r.HandleFunc(`/{key:[a-zA-Z][a-zA-Z\-\_]*}`, link).Methods("HEAD", "GET")
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

// Code copied from https://github.com/spf13/cobra.
func levenshteinDistance(s, t string) int {
	s = strings.ToLower(s)
	t = strings.ToLower(t)

	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}
