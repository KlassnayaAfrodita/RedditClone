package user

type User struct {
	ID       int    `json:"user_id"`
	Login    string `json:"user_login"`
	password string `json:"user_password`
}

type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	Register(login, pass string) (*User, error)
	GetUserByID(id int) (*User, error)
}
