package main

import (
	"fmt"
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

func linksHome() string {
	linksHome := os.Getenv("LINKS_HOME")
	if linksHome != "" {
		return linksHome
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, "links")
}

func redirect(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	fmt.Println("key", key)

	path := filepath.Join(home, key)
	fmt.Println("path", path)

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
	fmt.Println("url", url)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
