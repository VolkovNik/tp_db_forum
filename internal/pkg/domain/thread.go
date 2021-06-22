package domain

import "time"

type Thread struct {
	Id 		int 		`json:"id,omitempty"`
	Title 	string 		`json:"title,omitempty"`
	Author 	string 		`json:"author,omitempty"`
	Forum 	string 		`json:"forum,omitempty"`
	Message string 		`json:"message,omitempty"`
	Slug 	string 		`json:"slug,omitempty"`
	Created time.Time 	`json:"created,omitempty"`
	Votes 	int 		`json:"votes,omitempty"`
}

type Threads []Thread

type ThreadRepository interface {
	Create(thread *Thread) error
	GetThreadBySlug(slug string) (Thread, error)
	Vote(slugOrId string, vote Vote) (Thread, error)
	GetThreadBySlugOrId(slugOrId string) (Thread, error)
	UpdateThreadBySlugOrId(slugOrId string, thread *Thread) error
	GetForumThreads(slug string, limit int, since string, desc bool) (Threads, error)
}

type ThreadUsecase interface {
	Create(thread *Thread) error
	GetThreadBySlug(slug string) (Thread, error)
	Vote(slugOrId string, vote Vote) (Thread, error)
	GetThreadBySlugOrId(slugOrId string) (Thread, error)
	UpdateThreadBySlugOrId(slugOrId string, thread *Thread) error
	GetForumThreads(slug string, limit int, since string, desc bool) (Threads, error)
}