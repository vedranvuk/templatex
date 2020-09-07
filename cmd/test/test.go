package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/vedranvuk/templatex"
)

var nss *templatex.Namespaces

func handler(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	if p == "/" {
		p = "/home"
	}
	if err := nss.ExecuteNamespace(w, p, nil); err != nil {
		if errors.Is(err, templatex.ErrNamespaceNotFound) {
			nss.ExecuteNamespace(w, "/notfound", nil)
		}
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	}
}

func main() {
	log.SetFlags(0)

	ns, err := templatex.ParseRoot("../../test/testdata", "index", ".html")
	if err != nil {
		log.Fatal(err)
	}
	nss = ns

	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}
