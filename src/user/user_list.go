package user

var UserList []User = make([]User, 0, 1000)

func Join(user User) {
	UserList = append(UserList, user)
}
