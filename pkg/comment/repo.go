package comment

import (
	"errors"
	"sync"
)

var commentNotFound = errors.New("comment not found")

type CommentRepository struct {
	lastID int
	data   []*Comment
	mu     *sync.Mutex
}

func (repo *CommentRepository) GetAllByPost(postID int) ([]*Comment, error) {
	comments := make([]*Comment, 0, 10)
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, comment := range repo.data {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (repo *CommentRepository) GetAllByUser(userID int) ([]*Comment, error) {
	comments := make([]*Comment, 0, 10)
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, comment := range repo.data {
		if comment.UserID == userID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (repo *CommentRepository) GetByID(id int) (*Comment, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, comment := range repo.data {
		if comment.ID == id {
			return comment, nil
		}
	}
	return &Comment{}, commentNotFound
}

func (repo *CommentRepository) Add(comment *Comment) (int, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.lastID++
	comment.ID = repo.lastID
	repo.data = append(repo.data, comment)
	return comment.ID, nil
}

func (repo *CommentRepository) Delete(id int) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	i := -1
	for idx, comment := range repo.data {
		if comment.ID != id {
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

func (repo *CommentRepository) Update(newComment *Comment) (bool, error) {
	comment, err := repo.GetByID(newComment.id)
	if err != nil {
		return false, err
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()

	comment.Content = newComment.Content
	return true, nil
}
