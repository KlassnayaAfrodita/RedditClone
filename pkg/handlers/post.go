package handlers

import (
	"log/slog"
	"net/http"

	"github.com/KlassnayaAfrofita/RedditClone/pkg/post"
)

type PostHandler struct {
	Logger   *slog.Logger
	PostRepo post.PostRepo
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {

}
