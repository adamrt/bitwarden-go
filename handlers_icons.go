package main

import (
	"fmt"
	"net/http"
	"strings"
)

func iconHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) != 4 || p[1] != "icons" || p[3] != "icon.png" {
		http.NotFound(w, r)
		return
	}

	url := fmt.Sprintf("http://%s/favicon.ico", p[2])
	http.Redirect(w, r, url, 301)
}
