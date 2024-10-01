package post

import (
	"errors"
	"sync"
)

var productNotFound = errors.New("product not found")

type PostRepository struct {
	lastID int
	data   []*Post
	mu     *sync.Mutex
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		data: make([]*Post, 0, 10),
	}
}

func (repo *PostRepository) GetAll() ([]*Post, error) {
	return repo.data, nil
}

func (repo *PostRepository) GetByCategiry(category string) ([]*Post, error) {
	posts := make([]*Post, 0, 10)
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, post := range repo.data {
		if post.Category == category {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (repo *PostRepository) GetByUser(login string) ([]*Post, error) {
	posts := make([]*Post, 0, 10)
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, post := range repo.data {
		if post.CreatedBy == login {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (repo *PostRepository) GetByID(id int) (*Post, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, post := range repo.data {
		if post.ID == id {
			return post, nil
		}
	}
	return &Post{}, productNotFound
}

func (repo *PostRepository) Add(post *Post) (int, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.lastID++
	post.ID = repo.lastID
	repo.data = append(repo.data, post)
	return post.ID, nil
}

func (repo *PostRepository) Update(newPost *Post) (bool, error) {
	currentPost, err := repo.GetByID(newPost.ID)
	if err != nil {
		return false, err
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	currentPost.Title = newPost.Title
	currentPost.Content = newPost.Content
	currentPost.Category = newPost.Category
	currentPost.Rating = newPost.Rating
	return true, nil
}

func (repo *PostRepository) Delete(id int) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	i := -1
	for idx, post := range repo.data {
		if post.ID != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil // or the zero value of T
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}
