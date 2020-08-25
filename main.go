package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/matoous/go-nanoid"
	"net/http"
)

// make this type better
var links = make(map[string]string)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{identifier}", getLinkHandler).Methods("GET")
	r.HandleFunc("/links", createLinkHandler).Methods("POST")

	http.ListenAndServe(":8080", r)
}

type link struct {
	original   string
	identifier string
}

func createLink(original string) (*link, error) {
	// validate link
	link := &link{original: original}
	err := link.shorten()
	if err != nil {
		return link, err
	}

	links[link.identifier] = link.original

	return link, nil
}

func (link *link) shorten() error {
	identifier, err := gonanoid.ID(6)
	if err != nil {
		return err
	}

	link.identifier = identifier

	return nil
}

func (link *link) shortened(request *http.Request) string {
	return fmt.Sprintf("%s/%s", request.Host, link.identifier)
}

func createLinkHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	url := req.FormValue("url")
	if url == "" {
		http.Error(w, "Missing parameter: url", http.StatusBadRequest)
		return
	}

	link, err := createLink(url)
	if err != nil {
		http.Error(w, "Invalid parameter: url", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Your short url: %s\n", link.shortened(req))
}

func getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	identifier := vars["identifier"]
	url := links[identifier]

	if url == "" {
		msg := fmt.Sprintf("Unable to find %s\n", identifier)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	// log things

	http.Redirect(w, req, url, 302)
}
