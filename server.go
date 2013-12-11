
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
"crypto/rand"
"crypto/sha1"
"strconv"
"io"
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



type Game struct{

	GameName string
	GameID int
	HashValue string
	Salt string
	/*
	Timing string
	Minutes int
	Seconds int
	Increment int
	*/
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b % byte(len(alphanum))]
	}
	return string(bytes)
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

		gamename := r.FormValue("GameName");
		temppassword := r.FormValue("Password");
		salt := randString(20)

		h := sha1.New()
		io.WriteString(h, temppassword + salt)
		hashvalue := h.Sum(nil)

		/*
		timing := r.FormValue("Timing");
		minutes := r.FormValue("MinutesList");
		seconds := r.FormValue("SecondsList");
		increment := r.FormValue("IncrementList");
		*/

		stmt, err := db.Prepare("INSERT INTO Games(game_name, hash_value, salt) VALUES(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(gamename, hashvalue, salt)
		if err != nil {
			log.Fatal(err)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		val := strconv.FormatInt(lastId, 10)

		http.Redirect(w, r, "/game/"+val, http.StatusFound)


	default:

		var GameName string
		id, err := strconv.Atoi(arg)

		err = db.QueryRow("select game_name from Games where game_id = ?", id).Scan(&GameName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(GameName)

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

	var err error
	db, err = sql.Open("mysql",
		"root:password@tcp(127.0.0.1:3306)/ChessDatabase")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/game/", makeHandler(gameHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
