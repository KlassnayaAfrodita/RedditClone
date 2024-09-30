package session

import (
	"errors"
	"math/rand"
	"sync"
)

var sessionNotFound = errors.New("session not found")

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func createSessionToken() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type SessionRepository struct {
	data []*Session
	mu   *sync.Mutex
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		data: make([]*Session, 0, 10),
		mu:   &sync.Mutex{},
	}
}

func (repo *SessionRepository) GetUserID(sessionToken string) (*Session, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, session := range repo.data {
		if session.Token == sessionToken {
			return session, nil
		}
	}
	return &Session{}, sessionNotFound
}

func (repo *SessionRepository) Add(userID int) (string, error) {
	for _, session := range repo.data {
		if session.UserID == userID {
			_, err := repo.Update(userID)
			if err != nil {
				return "", err
			}
			return "", nil
		}
	}

	sessionToken := createSessionToken()

	for {
		_, err := repo.GetUserID(sessionToken)
		if err != nil {
			repo.mu.Lock()
			defer repo.mu.Unlock()
			repo.data = append(repo.data, &Session{
				UserID: userID,
				Token:  sessionToken,
			})
			return sessionToken, nil
		}
		sessionToken = createSessionToken()
	}
}

func (repo *SessionRepository) Update(userID int) (string, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, session := range repo.data {
		if session.UserID == userID {
			session = &Session{
				UserID: userID,
				Token:  createSessionToken(),
			}
			return session.Token, nil
		}
	}
	return "", sessionNotFound
}
