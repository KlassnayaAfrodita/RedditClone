package session

type Session struct {
	UserID int    `json:"user_id"`
	Token  string `json:"user_token"`
}

type SessionRepo interface {
	//TODO
	GetUserID(sessionToken string) (*Session, error)
	Add(userID int) (bool, error)
	Delete(userID int) (bool, error)
}
