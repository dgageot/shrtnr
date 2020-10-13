package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

var home = linksHome()

func main() {
	r := mux.NewRouter()
	r.HandleFunc(`/{key:[a-zA-Z][a-zA-Z\-\_]*}`, redirect).Methods("HEAD", "GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	path := filepath.Join(home, key)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	url := string(buf)
	url = strings.TrimSpace(url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func linksHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, "links")
}
