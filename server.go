
package main

import (
"html/template"
"io/ioutil"
"net/http"
"regexp"
"path/filepath"
)


const (
	tmplPath = "./templates"
	dataPath = "./data"
)

type Page struct {

	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title[len("data/"):len(title)-len(".txt")], Body: body}, nil
}

func landingHandler(w http.ResponseWriter, r *http.Request){

	var dataFile = "Chess.txt"
	var htmlTemplate = "Index.html" 

	data := filepath.Join(dataPath, dataFile)
	p, err := loadPage(data)
	if err != nil {
		p = &Page{Title: "Error"}
	}
	pagePath := filepath.Join(tmplPath, htmlTemplate)
	renderTemplate(w, p, pagePath)

}

func viewHandler(w http.ResponseWriter, r *http.Request, argument string){

	var dataFile string
	var htmlTemplate string

	switch argument{

	case "new":

		dataFile = "New Game.txt"
		htmlTemplate = "NewGame.html"  

	default:

		dataFile = "Error.txt"
		htmlTemplate = "Error.html" 

	}


	data := filepath.Join(dataPath, dataFile)
	p, err := loadPage(data)
	if err != nil {
		p = &Page{Title: "Error"}
	}
	pagePath := filepath.Join(tmplPath, htmlTemplate)
	renderTemplate(w, p, pagePath)

}



func renderTemplate(w http.ResponseWriter, p *Page,  path string) {
	var templates = template.Must(template.ParseFiles(path))
	err := templates.ExecuteTemplate(w, path[len("templates/"):], p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
