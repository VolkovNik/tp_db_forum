package domain

type Forum struct {
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title,omitempty"`
	User 	string `json:"user,omitempty"`
	Posts   int    `json:"posts,omitempty"`
	Threads int    `json:"threads,omitempty"`
}

type ForumRepository interface {
	Create(forum *Forum) (*Forum, error)
	GetForumDetails(slug string) (*Forum, error)
	GetForumUsers(slug string, limit int, since string, desc bool) (Users, error)
}

type ForumUsecase interface {
	Create(forum *Forum) (*Forum, error)
	GetForumDetails(slug string) (*Forum, error)
	GetForumUsers(slug string, limit int, since string, desc bool) (Users, error)
}