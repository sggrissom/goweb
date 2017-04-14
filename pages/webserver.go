package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"time"
	"io/ioutil"
)

var pages []page
var templates = template.Must(template.ParseGlob("templates/*/*.tmpl"))

func getPageObject(path string) page {
	for _, page := range pages {
		if page.Path == path {
			return page
		}
	}

	return page{Path: "notFound", Title: "404 not found"}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var path string

	if len(r.URL.Path) > 1 {
		path = r.URL.Path[1:]
	} else {
		path = "home"
	}
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

	Page := getPageObject(path)

	if len(Page.Title) == 0 {
		Page.Title = "Go Website"
	}

	var pageData map[string]interface{}
	if Page.DataFn != nil {
		pageData = Page.DataFn()
	}

	data := struct {
		Page     string
		Title    string
		Style    string
		JS       []string
		Nav      []page
		PageData map[string]interface{}
	}{
		Page.Path,
		Page.Title,
		style,
		Page.JS,
		pages,
		pageData,
	}

	if len(Page.Template) > 0 {
		path = Page.Template
	}

	renderTemplate(w, path, data)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

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

func query(db *sql.DB) {
	var (
		id    int
		first string
		last  string
	)
	rows, err := db.Query("select * from players")
	if err != nil {
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &first, &last)
		if err != nil {
		}

		fmt.Printf("%d. %s %s", id, first, last)
	}
	err = rows.Err()
	if err != nil {
	}
}

func main() {

	password, err := ioutil.ReadFile("password.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	password = password[:len(password)-1]
	
	fmt.Printf("password: %s\n", password)

	dbConnectionString := fmt.Sprintf("website:%s@tcp(steven.db:3306)/ball", password)

	fmt.Printf("connectionString: %s\n", dbConnectionString)
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	defer db.Close()

		fmt.Printf("before ping\n")
	err = db.Ping()
		fmt.Printf("pinged\n")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
		fmt.Printf("success!!!!\n")

	query(db)

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
