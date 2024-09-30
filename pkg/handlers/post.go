package handlers

import "log/slog"

type PostHandler struct {
	Logger   *slog.Logger
	PostRepo post.PostRepo
}
