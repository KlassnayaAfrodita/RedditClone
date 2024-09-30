package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/KlassnayaAfrodita/RedditClone/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/login": struct{}{},
	}
	noSessUrls = map[string]struct{}{
		"/": struct{}{},
	}
)

// TODO поставить последним
func AuthorizeMiddleware(logger slog.Logger, sr *session.SessionRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Resource: AuthorizeMiddleware")

		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session_id")
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			http.Error(w, "you dont auth", 200)
			http.Redirect(w, r, "/login", 200)
			return
		}

		sessionUser := cookie.Value
		sess, err := sr.GetUserID(sessionUser)
		if err != nil {
			http.Error(w, "db error", 500)
			logger.Error("Resource: AuthorizeMiddleware",
				"Error",
				err,
			)
			return
		}

		type sessionKey string
		var SessionKey sessionKey = "session_id"
		sess, err = sr.Add(sess.UserID)
		ctx := context.WithValue(r.Context(), SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
