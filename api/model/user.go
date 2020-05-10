package model

import (
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
		log.Fatalf(err.Error())
	}
	defer rows.Close()
	var result = make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
			log.Fatalf(err.Error())
		}
		result = append(result, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf(err.Error())
	}
	return result
}

/*
	find user by id
*/
func (_ User) FindById(id int) *User {
	for _, user := range userRepo {
		if user.Id == id {
			return &user
		}
	}
	return nil
}

/*
	save new user in a repository
*/
func (user *User) Save() {
	user.Id = len(userRepo) + 1
	userRepo = append(userRepo, *user)
}

/*
	delete user in a repository
*/
func (user User) Remove() {
	for key, u := range userRepo {
		if u.Id == user.Id {
			leftSlice := userRepo[0:key]
			rightSlice := userRepo[key+1:]
			userRepo = append(leftSlice, rightSlice...)
			break
		}
	}
}

func (user *User) Update(id int) {
	index := -1
	for key, u := range userRepo {
		if u.Id == id {
			index = key
			break
		}
	}
	u := &userRepo[index]
	u.Id = id
	u.Name = user.Name
	u.Age = user.Age
}
