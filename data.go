package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid     string            `valid:"required,uuidv4"`
	Username string            `valid:"required,alphanum"`
	Password string            `valid:"required"`
	Fname    string            `valid:"required,alpha"`
	Lname    string            `valid:"required,alpha"`
	Email    string            `valid:"required,email"`
	Errors   map[string]string `valid:"-"`
}

func saveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (uuid text, firstname text, lastname text, username text, email text, password text)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (uuid, firstname, lastname, username, email, password) values (?,?,?,?,?,?)")
	_, err := stmt.Exec(u.Uuid, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

func userExists(u *User) bool {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	var ps, us string
	q, err := db.Query("select username, password from users where username ='" + u.Username + "'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&us, &ps)
	}

	pw := bcrypt.CompareHashAndPassword([]byte(ps), []byte(u.Password))
	if us == u.Username && pw == nil {
		return true
	}
	return false
}

func encryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

func Uuid() string {
	id := uuid.Must(uuid.NewV4())
	return id.String()
}
