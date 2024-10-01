package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KlassnayaAfrodita/RedditClone/pkg/comment"
	"github.com/KlassnayaAfrodita/RedditClone/pkg/session"
	"github.com/KlassnayaAfrofita/RedditClone/pkg/post"
	"github.com/KlassnayaAfrofita/RedditClone/pkg/user"
	"github.com/gorilla/mux"
)

type sessionKey string

var SessionKey sessionKey = "session_id"

type PostHandler struct {
	Logger      *slog.Logger
	PostRepo    post.PostRepo
	UserRepo    user.UserRepo
	CommentRepo comment.CommentRepo
	SessionRepo session.SessionRepo
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostRepo.GetAll()
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.GetAll",
			"error:", err)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, `{"error": "json error"}`, 500)
		h.Logger.Error("Resource: handler.post.GetAll",
			"error:", err)
		return
	}

	w.Write(resp)
}

func (h *PostHandler) AddPost(w http.ResponseWriter, r *http.Request) { //* принимаем post запрос с json
	// ctx := r.Context()
	// if ctx != nil {
	// 	http.Error(w, `{"error": "context error"}`, 400)
	// 	h.Logger.Error("Resource: handler.post.AddPost",
	// 		"error:", ctx)
	// 	return
	// }
	// userSession := ctx.Value(SessionKey).(*session.Session)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.Logger.Error("Resource: handler.post.AddPost",
			"error: wrong metgod")
		return
	}

	var post *post.Post

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, `{"error": "input error"}`, 400)
		h.Logger.Error("Resource: handler.post.AddPost",
			"error:", err)
		return
	}

	err = json.Unmarshal(body, post)
	if err != nil {
		http.Error(w, `{"error": "json error"}`, 500)
		h.Logger.Error("Resource: handler.post.AddPost",
			"error:", err)
		return
	}

	ctx := r.Context()
	user := ctx.Value("session_id").(*session.Session)
	us, err := h.UserRepo.GetUserByID(user.UserID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.AddPost",
			"error:", err)
		return
	}
	post.CreatedBy = us.Login

	_, err = h.PostRepo.Add(post)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.AddPost",
			"error:", err)
		return
	}

	http.Redirect(w, r, "/posts", 200)
}

func (h *PostHandler) GetByCategory(w http.ResponseWriter, r *http.Request) { //* получаем имя категории через url
	vars := mux.Vars(r)
	category := vars["CATEGORY_NAME"]

	posts, err := h.PostRepo.GetByCategiry(category)
	if err != nil {
		h.Logger.Error("Resource: handler.post.GetByCategory",
			"error:", err)
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		h.Logger.Error("Resource: handler.post.GetByCategory",
			"error:", err)
		http.Error(w, `{"error": "json error"}`, 500)
		return
	}

	w.Write(resp)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) { //* получаем id поста через url
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil {
		h.Logger.Error("Resource: handler.post.GetPost",
			"error:", err)
		http.Error(w, `{"error": "input error"}`, 400)
		return
	}

	post, err := h.PostRepo.GetByID(id)
	if err != nil {
		h.Logger.Error("Resource: handler.post.GetPost",
			"error:", err)
		http.Error(w, `{"error": "db error"}`, 500)
		return
	}

	resp, err := json.Marshal(post)
	if err != nil {
		h.Logger.Error("Resource: handler.post.GetPost",
			"error:", err)
		http.Error(w, `{"error": "json error"}`, 500)
		return
	}

	w.Write(resp)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) { //* получаем через post json
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.Logger.Error("Resource: handler.post.AddComment",
			"error: wrong metgod")
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.Logger.Error("Resource: handler.post.AddComment",
			"error:", err)
		http.Error(w, `{"error": "input error"}`, 400)
		return
	}

	var comment *comment.Comment

	err = json.Unmarshal(body, comment)
	if err != nil {
		h.Logger.Error("Resource: handler.post.AddComment",
			"error:", err)
		http.Error(w, `{"error": "json error"}`, 500)
		return
	}

	ctx := r.Context()
	user := ctx.Value("session_id").(*session.Session)
	us, err := h.UserRepo.GetUserByID(user.UserID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.AddComment",
			"error:", err)
		return
	}
	comment.CreatedBy = us.Login //! проставляем автора коммента в поле

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	comment.PostID = id //! проставляем пост, к которому привязан коммент

	_, err = h.PostRepo.Add(comment)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.AddComment",
			"error:", err)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(id), 200)
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) { //* метод Delete, в url получаем id поста
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil || r.Method != http.MethodDelete {
		http.Error(w, `{"error": "input error"}`, 400)
		h.Logger.Error("Resource: handler.post.DeleteComment",
			"error:", err)
		return
	}

	_, err = h.CommentRepo.Delete(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.DeleteComment",
			"error:", err)
		return
	}

	http.Redirect(w, r, "/posts", 200)
}

func (h *PostHandler) PostUpvote(w http.ResponseWriter, r *http.Request) { //* получаем id из url
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil {
		http.Error(w, `{"error": "input error"}`, 400)
		h.Logger.Error("Resource: handler.post.PostUpvote",
			"error:", err)
		return
	}

	post, err := h.PostRepo.GetByID(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.PostUpvote",
			"error:", err)
		return
	}

	post.Rating++

	http.Redirect(w, r, "/post/"+strconv.Itoa(id), 200)
}

func (h *PostHandler) PostDownvote(w http.ResponseWriter, r *http.Request) { //* получаем id из url
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil {
		http.Error(w, `{"error": "input error"}`, 400)
		h.Logger.Error("Resource: handler.post.PostDownvote",
			"error:", err)
		return
	}

	post, err := h.PostRepo.GetByID(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.PostDownvote",
			"error:", err)
		return
	}

	post.Rating--

	http.Redirect(w, r, "/post/"+strconv.Itoa(id), 200)
}

// TODO
func (h *PostHandler) PostUnvote(w http.ResponseWriter, r *http.Request) { //* получаем id из url

}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) { //* метод Delete, id получаем из url
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["POST_ID"])
	if err != nil || r.Method != http.MethodDelete {
		http.Error(w, `{"error": "input error"}`, 400)
		h.Logger.Error("Resource: handler.post.DeletePost",
			"error:", err)
		return
	}

	_, err = h.PostRepo.Delete(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.DeletePost",
			"error:", err)
		return
	}

	http.Redirect(w, r, "/posts", 200)
}

func (h *PostHandler) GetByUser(w http.ResponseWriter, r *http.Request) { //* получаем userID из url
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["USER_ID"])
	if err != nil {
		http.Error(w, `{"error": "input error"}`, 400)
		h.Logger.Error("Resource: handler.post.GetByUser",
			"error:", err)
		return
	}

	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.GetByUser",
			"error:", err)
		return
	}
	posts, err := h.PostRepo.GetByUser(user.Login)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, 500)
		h.Logger.Error("Resource: handler.post.GetByUser",
			"error:", err)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, `{"error": "json error"}`, 500)
		h.Logger.Error("Resource: handler.post.GetByUser",
			"error:", err)
		return
	}

	w.Write(resp)
}
