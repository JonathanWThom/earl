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

func createLink(original string) *link {
	link := &link{original: original}
	link.shorten()
	links[link.identifier] = link.original

	return link
}

func (link *link) shorten() {
	identifier, err := gonanoid.ID(6)
	if err != nil {
		panic(err)
	}

	link.identifier = identifier
}

func (link *link) shortened(request *http.Request) string {
	return fmt.Sprintf("%s/%s", request.Host, link.identifier)
}

func createLinkHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	url := req.FormValue("url")
	// handle error
	link := createLink(url)

	fmt.Fprintf(w, "Your short url: %s\n", link.shortened(req))
}

func getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Println(links)
	url := links[vars["identifier"]]

	// log things
	// redirect here
	// handle empty

	http.Redirect(w, req, url, 302)
	//fmt.Fprintf(w, "%s\n", url)
}
