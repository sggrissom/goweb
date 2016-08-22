package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pages = make(map[string]func(string))

var templates = template.Must(template.ParseGlob("templates/*/*.tmpl"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", p)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateCookie(w http.ResponseWriter, Name string, Value string) {
	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{
		Name:     Name,
		Value:    Value,
		Expires:  expiration,
		HttpOnly: false,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path[1:]
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	style := r.FormValue("button")

	if len(style) > 0 {
		updateCookie(w, "Style", style)
	} else {
		cookie, _ := r.Cookie("Style")

		if cookie == nil {
			style = "flatly"
			updateCookie(w, "Style", style)
		} else {
			style = cookie.Value
		}
	}

	data := struct {
		Page  string
		Title string
		Style string
	}{
		path,
		"Go Website",
		style,
	}

	for pageName, pageFunc := range pages {
		if pageName == path {
			pageFunc(pageName)
		}
	}

	renderTemplate(w, path, data)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func main() {
	RegisterPages()

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/static/favicon.ico")
	})
	http.HandleFunc("/", makeHandler(handler))
	http.ListenAndServe(":8080", nil)
}
