package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

type Article struct {
	Id                    uint
	Title, Anon, FullText string
}

var news []Article
var showA = Article{}

	func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err)
	}
	news = []Article{}
	for result.Next() {
		var a Article
		err = result.Scan(&a.Id, &a.Title, &a.Anon, &a.FullText)
		if err != nil {
			panic(err)
		}
		news = append(news, a)
	}

	t.ExecuteTemplate(w, "index", news)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anon := r.FormValue("anon")
	fullText := r.FormValue("full-text")

	if title == "" || anon == "" || fullText == "" {
		fmt.Fprint(w, "Any fields are not filled in")
	} else {

		db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3307)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO articles (title,anon,full_text) VALUES (?,?,?)", title, anon, fullText)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func showArticle(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show_article.html", "templates/footer.html", "templates/header.html")
	if err != nil{
		fmt.Fprintf(w, err.Error())
	}


	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id=%s", vars["id"]))
	if err != nil {
		panic(err)
	}

	showA = Article{}
	for result.Next() {
		var a Article
		err = result.Scan(&a.Id, &a.Title, &a.Anon, &a.FullText)
		if err != nil {
			panic(err)
		}
		showA = a
	}

	t.ExecuteTemplate(w, "show_article", showA)
}

func signIn(w http.ResponseWriter, r *http.Request)  {
	t, err := template.ParseFiles("templates/sign_in.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "signin", nil)
}

func deleteArticle(w http.ResponseWriter, r *http.Request)  {
	s := r.Referer()

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	a, err := strconv.ParseInt(string(s[len(s)- 1]), 10, 64)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM articles WHERE id = ?", a)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/save_article", saveArticle).Methods("POST")
	router.HandleFunc("/article/{id:[0-9]+}", showArticle).Methods("GET")
	router.HandleFunc("/sign_in", signIn)
	router.HandleFunc("/delete_article", deleteArticle).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

func main() {
	handleFunc()
}
