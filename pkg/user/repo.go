package user

import (
	"errors"
	"math/rand"
	"sync"
)

var (
	ErrNoUser     = errors.New("No user found")
	ErrBadPass    = errors.New("Invald password")
	ErrUserExists = errors.New("User already exists")
)

func RandID() int {
	return rand.Int()
}

type UserRepository struct {
	lastID int
	mu     *sync.Mutex
	data   map[string]*User //! логин - юзер
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		mu: &sync.Mutex{},
		data: map[string]*User{
			"eyenot": &User{
				ID:       1,
				Login:    "eyenot",
				password: "123",
			},
		},
	}
}

func (repo *UserRepository) Authorize(login, pass string) (*User, error) {
	user, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}
	if user.password != pass {
		return nil, ErrBadPass
	}
	return user, nil
}

func (repo *UserRepository) Register(login, pass string) (*User, error) {
	if user, ok := repo.data[login]; ok {
		return user, ErrUserExists
	}
	newID := RandID()
	for _, user := range repo.data {
		if newID == user.ID {
			newID = RandID()
		}
	}
	return &User{ID: newID, Login: login, password: pass}, nil
}

func (repo *UserRepository) Register(login, pass string) (*User, error) {
	for _, user := range repo.data {
		if user.Login == login {
			return user, ErrUserExists
		}
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.lastID++
	newUser := &User{
		ID:       repo.lastID,
		Login:    login,
		password: pass,
	}
	return newUser, nil
}

func (repo *UserRepository) GetUserByID(id int) (*User, error) {
	for _, user := range repo.data {
		if user.ID == id {
			return user, nil
		}
	}
	return &User{}, ErrNoUser
}
