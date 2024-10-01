package comment

type Comment struct {
	ID        int    `json:"comment_id,omitempty"`
	PostID    int    `json:"comment_post_id"`
	Content   string `json:"comment_content"`
	CreatedBy string `json:"comment_created_by"`
}

type CommentRepo interface {
	GetAllByPost(postID int) ([]*Comment, error)
	GetAllByUser(createdBy string) ([]*Comment, error)
	GetByID(id int) (*Comment, error)
	Add(comment *Comment) (int, error)
	Delete(id int) (bool, error)
	Update(newComment *Comment) (bool, error)
}
