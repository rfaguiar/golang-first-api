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

/*
	find user by id
*/
func (_ User) FindById(id int) *User {
	for _, user := range UserRepo {
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
	user.Id = len(UserRepo) + 1
	UserRepo = append(UserRepo, *user)
}

/*
	delete user in a repository
*/
func (user User) Remove() {
	for key, u := range UserRepo {
		if u.Id == user.Id {
			leftSlice := UserRepo[0:key]
			rightSlice := UserRepo[key+1:]
			UserRepo = append(leftSlice, rightSlice...)
			break
		}
	}
}

func (user *User) Update(id int) {
	index := -1
	for key, u := range UserRepo {
		if u.Id == id {
			index = key
			break
		}
	}
	u := &UserRepo[index]
	u.Id = id
	u.Name = user.Name
	u.Age = user.Age
}
