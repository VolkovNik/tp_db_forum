package domain

import "time"

type Post struct {
	Id 			int 		`json:"id,omitempty"`
	Parent		int 		`json:"parent,omitempty"`
	Author 		string 		`json:"author,omitempty"`
	Message 	string 		`json:"message,omitempty"`
	IsEdited 	bool		`json:"isEdited,omitempty"`
	Forum 		string 		`json:"forum,omitempty"`
	Thread 		int 		`json:"thread,omitempty"`
	Created 	time.Time 	`json:"created,omitempty"`
}

type PostFull struct {
	Post 	*Post 	`json:"post,omitempty"`
	Author 	*User 	`json:"author,omitempty"`
	Thread	*Thread `json:"thread,omitempty"`
	Forum 	*Forum 	`json:"forum,omitempty"`
}

type Posts []Post

func (p Posts) Size() int{
	return len(p)
}

type PostRepository interface {
	Create(slugOrId string, posts *Posts) error
	UpdateById(id int, post *Post) error
	Details(id int, related string) (PostFull, error)
	GetThreadPosts(slugOrId string, limit int, since int, sort string, desc bool) (Posts, error)
}

type PostUsecase interface {
	Create(slugOrId string, posts *Posts) error
	UpdateById(id int, post *Post) error
	Details(id int, related string) (PostFull, error)
	GetThreadPosts(slugOrId string, limit int, since int, sort string, desc bool) (Posts, error)
}