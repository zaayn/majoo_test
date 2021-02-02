package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Nama string
	Email string
	Foto string
	Password string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM User ORDER BY id ASC")
    if err != nil {
        panic(err.Error())
    }
    usr := User{}
    res := []User{}
    for selDB.Next() {
        var id int
		var name string
		var email string
		var foto string
		var password string
        err = selDB.Scan(&id, &name, &email, &foto, &password)
        if err != nil {
            panic(err.Error())
        }
        usr.ID = id
        usr.Nama = name
		usr.Email = email
		usr.Foto = foto
		usr.Password = password
        res = append(res, usr)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM User WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
		
	}
	usr := User{}
	for selDB.Next() {
		var id int
		var name, email, foto, password string
		err = selDB.Scan(&id, &name, &email, &foto, &password)
		if err != nil {
			panic(err.Error())
		}
		usr.ID = id
		usr.Nama = name
		usr.Email = email
		usr.Foto = foto
		usr.Password = password
	}
	tmpl.ExecuteTemplate(w, "Show", usr)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM User WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	usr := User{}
	for selDB.Next() {
		var id int
		var name, email, foto, password string
		err = selDB.Scan(&id, &name, &email, &foto, &password)
		if err != nil {
			panic(err.Error())
		}
		usr.ID = id
		usr.Nama = name
		usr.Email = email
		usr.Foto = foto
		usr.Password = password
	}
	tmpl.ExecuteTemplate(w, "Edit", usr)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		foto := r.FormValue("foto")
		password := r.FormValue("password")
		insForm, err := db.Prepare("INSERT INTO User(name, email, foto, password) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, foto, password)
		// log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		foto := r.FormValue("foto")
		password := r.FormValue("password")
		insForm, err := db.Prepare("UPDATE User SET name=?, email=?, foto=?, password=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, foto, password)

	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	usr := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM User WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(usr)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {

	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
