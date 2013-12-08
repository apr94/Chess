
package main

import (
"html/template"
"io/ioutil"
"net/http"
"regexp"
"path/filepath"
"database/sql"
_ "github.com/go-sql-driver/mysql"
"log"
"fmt"
)


const (
	tmplPath = "./templates"
	dataPath = "./data"
)

var db *sql.DB

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

func viewHandler(w http.ResponseWriter, r *http.Request, arg string){

	var dataFile string
	var htmlTemplate string

	switch arg{

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

func gameHandler(w http.ResponseWriter, r * http.Request, arg string){


	switch arg{

	case "":

		stmt, err := db.Prepare("select game_id, game_name from Games where game_id = 1")
		if err != nil {
			log.Fatal(err)
		}
		var name string
		err = stmt.QueryRow(1).Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)


	default:

	}



}



func renderTemplate(w http.ResponseWriter, p *Page,  path string) {
	var templates = template.Must(template.ParseFiles(path))
	err := templates.ExecuteTemplate(w, path[len("templates/"):], p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(view|game)/([a-zA-Z0-9]*)$")

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
	http.HandleFunc("/game/", makeHandler(gameHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)

	db, err := sql.Open("mysql",
		"root:password@tcp(127.0.0.1:3306)/ChessDatabase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
