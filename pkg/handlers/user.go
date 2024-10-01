package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/KlassnayaAfrodita/RedditClone/pkg/session"

	"github.com/KlassnayaAfrodita/RedditClone/pkg/user"
)

// type sessionKey string

// var SessionKey sessionKey = "session_id"

type UserHandler struct {
	Logger      *slog.Logger
	UserRepo    *user.UserRepo
	SessionRepo *session.SessionRepo
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) { //* получаем из post json
	if r.Method == http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.Logger.Error("Resource: handler.user.Login",
			"error: wrong metgod")
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, `{"error": "input error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	var user *user.User

	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, `{"error": "json error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	user, err = h.UserRepo.Authorize(user.Login, user.Password)
	if err != nil {
		http.Error(w, `{"error": "auth error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	sess, err := h.SessionRepo.Update(user.ID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	h.Logger.Info("Resource: handler.user.Login", "user auth")

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/post", 200)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) { //* получаем из post json
	if r.Method == http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.Logger.Error("Resource: handler.user.Login",
			"error: wrong metgod")
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, `{"error": "input error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	var user *user.User

	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, `{"error": "json error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	user, err = h.UserRepo.Register(user.Login, user.Password)
	if err != nil {
		http.Error(w, `{"error": "auth error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	sess, err := h.SessionRepo.Add(user.ID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.user.Login",
			"error:", err)
		return
	}

	h.Logger.Info("Resource: handler.user.Login", "user auth")

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess.Token,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/post", 200)
}

// TODO удалить контекст, сделать недействительными куки, редирект
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sess := ctx.Value(SessionKey).(*session.Session)

	newCtx := context.WithValue(ctx, SessionKey, *session.Session)

	cookie := &http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
	}

	h.SessionRepo.Delete(sess.UserID)

	http.SetCookie(w, cookie)
	http.Redirect(w, r.WithContext(newCtx), "/", 200)
}
