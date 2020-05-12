package model

import (
	"database/sql"
	"fmt"
	"github.com/rfaguiar/golang-first-api/database"
	"log"
)

/*
	User type for use CAD
*/
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/*
	Repository for users type
*/
var userRepo = make([]User, 0)

/*
	find all users in a repository
*/
func (_ User) FindAll() []User {
	db := database.Current()
	rows, err := db.Query("select id, name, age from person")
	if err != nil {
		log.Print(err.Error())
	}
	defer rows.Close()
	var result = make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
			log.Print(err.Error())
		}
		result = append(result, u)
	}
	if err := rows.Err(); err != nil {
		log.Print(err.Error())
	}
	return result
}

/*
	find user by id
*/
func (_ User) FindById(id int) *User {
	var u User
	db := database.Current()
	row := db.QueryRow("select id, name, age from person where id = $1", id)
	if err := row.Scan(&u.Id, &u.Name, &u.Age); err != nil {
		log.Print(err.Error())
		return nil
	}
	return &u
}

/*
	save new user in a repository
*/
func (user *User) Save() error {
	errStr := fmt.Sprintf("Error when create user %v", user)
	tx := getDbTrancation()
	stmt, err := tx.Prepare("insert into person (name, age) values ($1, $2)")
	if err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(user.Name, user.Age); err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	return nil
}

/*
	delete user in a repository
*/
func (user User) Remove() error {
	errStr := fmt.Sprintf("Error when delete user %v", user)
	tx := getDbTrancation()
	stmt, err := tx.Prepare("delete from person where id = $1")
	if err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(user.Id); err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	return nil
}

/*
	update user in a repository
*/
func (user *User) Update(id int) error {
	errStr := fmt.Sprintf("Error when update user %v", user)
	tx := getDbTrancation()
	stmt, err := tx.Prepare("update person set name = $2, age = $3 where id = $1")
	if err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(id, user.Name, user.Age); err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Print(err.Error())
		return fmt.Errorf(errStr)
	}
	return nil
}

func getDbTrancation() *sql.Tx {
	tx, err := database.Current().Begin()
	if err != nil {
		log.Print(err.Error())
	}
	return tx
}
