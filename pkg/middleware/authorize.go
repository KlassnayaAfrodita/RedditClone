package middleware

import (
	"context"
	"log/slog"
	"net/http"
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

		session, err := r.Cookie("session_id")
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			http.Error(w, "you dont auth", 200)
			http.Redirect(w, r, "/login", 200)
			return
		}

		session, err = sr.GetUserID(session)
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
		session, err = sr.Add(session.UserID)
		ctx := context.WithValue(r.Context(), SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
