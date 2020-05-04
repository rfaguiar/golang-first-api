package model

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
var UserRepo []User

/*
	Initialize in memory database
*/
func (_ User) InitializeDatabase() {
	UserRepo = append(UserRepo,
		User{Id: 1, Name: "Jorge", Age: 20},
		User{Id: 2, Name: "Jhon", Age: 33})
}

/*
	find all users in a repository
*/
func (_ User) FindAll() []User {
	return UserRepo
}
