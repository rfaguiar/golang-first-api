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
	tx := getDbTransaction()
	stmt, err := tx.Prepare("insert into person (name, age) values ($1, $2)")
	if err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(user.Name, user.Age); err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	row := database.Current().QueryRow("select id from person where name = $1 and age = $2", user.Name, user.Age)
	if err := row.Scan(&user.Id); err != nil {
		log.Printf("Error when find user by name and age: %s", err.Error())
	}
	return nil
}

/*
	delete user in a repository
*/
func (user User) Remove() error {
	tx := getDbTransaction()
	stmt, err := tx.Prepare("delete from person where id = $1")
	if err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(user.Id); err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	return nil
}

/*
	update user in a repository
*/
func (user *User) Update(id int) error {
	tx := getDbTransaction()
	stmt, err := tx.Prepare("update person set name = $2, age = $3 where id = $1")
	if err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(id, user.Name, user.Age); err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return rollbackTransactionLogError(tx, err)
	}
	return nil
}

/*
	Transaction rollback.
	Log error.
	Return transaction error formated
*/
func rollbackTransactionLogError(tx *sql.Tx, err error) error {
	tx.Rollback()
	log.Print(err.Error())
	return fmt.Errorf("Error when execute transcation")
}

/*
	Create database unique secure transaction
	If error log in a console
*/
func getDbTransaction() *sql.Tx {
	tx, err := database.Current().Begin()
	if err != nil {
		log.Print(err.Error())
	}
	return tx
}
