package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

var router = mux.NewRouter()

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/example", examplePage)
	router.HandleFunc("/signup", signup).Methods("POST", "GET")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("signup.html", "index.html", "base.html")
		u := &User{}
		tmpl.ExecuteTemplate(w, "base", u)
	case "POST":
		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("email")
		un := r.FormValue("userName")
		pass := r.FormValue("password")

		u := &User{Fname: f, Lname: l, Email: em, Username: un, Password: pass}
		saveData(u)
		http.Redirect(w, r, "/", 302)
	}
}

func examplePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("base.html", "index.html", "internal.html")
	username := getUsername(r)
	if username != "" {
		err := tmpl.ExecuteTemplate(w, "base", &User{Username: username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first"))
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	pass := r.FormValue("password")
	u := &User{Username: name, Password: pass}

	redirect := "/"
	if name != "" && pass != "" {
		if userExists(u) {
			setSession(&User{Username: name, Password: pass}, w)
			redirect = "/example"
		} else {
			setMsg(w, "message", []byte("Please, signup or enter a valid username and password"))
		}
	} else {
		setMsg(w, "message", []byte("Username or password field are empty!"))
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	msg, _ := getMsg(w, r, "message")
	if msg != nil {
		tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html", "flash.html")
		err := tmpl.ExecuteTemplate(w, "base", msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		u := &User{}
		tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html")
		err := tmpl.ExecuteTemplate(w, "base", u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
