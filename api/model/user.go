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
var userRepo = make([]User, 0)

/*
	find all users in a repository
*/
func (_ User) FindAll() []User {
	return userRepo
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
