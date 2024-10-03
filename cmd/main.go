package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/KlassnayaAfrodita/RedditClone/pkg/comment"
	"github.com/KlassnayaAfrodita/RedditClone/pkg/handlers"
	"github.com/KlassnayaAfrodita/RedditClone/pkg/middleware"
	"github.com/KlassnayaAfrodita/RedditClone/pkg/post"
	"github.com/KlassnayaAfrodita/RedditClone/pkg/session"
	"github.com/KlassnayaAfrodita/RedditClone/pkg/user"
	"github.com/gorilla/mux"
)

func main() {
	sr := session.NewSessionRepository()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userRepo := user.NewUserRepository()
	postRepo := post.NewPostRepository()
	commentRepo := comment.NewCommentRepository()
	sessionRepo := session.NewSessionRepository()

	userHandler := &handlers.UserHandler{
		Logger:      logger,
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}

	postHandler := &handlers.PostHandler{
		Logger:      logger,
		PostRepo:    postRepo,
		UserRepo:    userRepo,
		CommentRepo: commentRepo,
		SessionRepo: sessionRepo,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/api/logout", userHandler.Logout).Methods("GET")

	router.HandleFunc("/api/posts/", postHandler.GetAll).Methods("GET")
	router.HandleFunc("/api/posts/", postHandler.AddPost).Methods("POST")
	router.HandleFunc("/api/posts/{CATEGORY_NAME}", postHandler.GetByCategory).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}", postHandler.GetPost).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}", postHandler.AddComment).Methods("POST")
	router.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postHandler.DeleteComment()).Methods("DELETE")
	router.HandleFunc("/api/post/{POST_ID}/upvote", postHandler.PostUpvote).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}/downvote", postHandler.PostDownvote).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}/unvote", postHandler.PostUnvote).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}", postHandler.DeletePost).Methods("DELETE")
	router.HandleFunc("/api/user/{USER_LOGIN}", postHandler.GetByUser).Methods("GET")

	mux := middleware.Auth(sr, router)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Info("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
