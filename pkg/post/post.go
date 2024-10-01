package post

type Post struct {
	ID        int    `json:"post_id,omitempty,omitempty"`
	Title     string `json:"post_title"`
	Content   string `json:"post_body"`
	CreatedBy string `json:"post_created_by,omitempty"`
	Category  string `json:"post_category"`
	Rating    int    `json:"post_rating"`
}

type PostRepo interface {
	GetAll() ([]*Post, error)
	GetByCategiry(category string) ([]*Post, error)
	GetByUser(login string) ([]*Post, error)
	GetByID(id int) (*Post, error)
	Add(post *Post) (int, error)
	Update(newPost *Post) (bool, error)
	Delete(id int) (bool, error)
}
