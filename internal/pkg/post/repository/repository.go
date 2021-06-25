package repository

import (
	"TPForum/internal/pkg/domain"
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
	"strings"
)

type PostRepository struct {
	db *pgx.ConnPool
}

var (
	selectFlatSinceDesc   = "SELECT id, author, forum, thread, message, parent, isEdited, created FROM posts WHERE thread = $1 AND id < $2 ORDER BY id DESC LIMIT $3"
	selectFlatSinceAsc 	  = "SELECT id, author, forum, thread, message, parent, isEdited, created FROM posts WHERE thread = $1 AND id > $2 ORDER BY id ASC LIMIT $3"

	selectTreeSinceDesc   = "SELECT posts.id, posts.author, posts.forum, posts.thread, posts.message, posts.parent, posts.isEdited, posts.created FROM posts JOIN posts P ON P.id = $2 WHERE posts.path < p.path AND posts.thread = $1 ORDER BY posts.path[1] DESC, posts.path DESC LIMIT $3"
	selectTreeSinceAsc    = "SELECT posts.id, posts.author, posts.forum, posts.thread, posts.message, posts.parent, posts.isEdited, posts.created FROM posts JOIN posts P ON P.id = $2 WHERE posts.path > p.path AND posts.thread = $1 ORDER BY posts.path[1] ASC, posts.path ASC LIMIT $3"

	selectParentSinceDesc = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created FROM posts as p WHERE p.thread = $1 AND p.path && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 AND p.path < (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ORDER BY p.path[1] DESC, p.path LIMIT $3)) ORDER BY p.path[1] DESC, p.path"
	selectParentSinceAsc  = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created FROM posts as p WHERE p.thread = $1 AND p.path && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 AND p.path > (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ORDER BY p.path[1] ASC, p.path LIMIT $3)) ORDER BY p.path[1] ASC, p.path"

	selectTreeDesc 		  = "SELECT posts.id, posts.author, posts.forum, posts.thread, posts.message, posts.parent, posts.isEdited, posts.created FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] DESC, posts.path DESC LIMIT $2"
	selectTreeAsc		  = "SELECT posts.id, posts.author, posts.forum, posts.thread, posts.message, posts.parent, posts.isEdited, posts.created FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] ASC, posts.path ASC LIMIT $2"

	selectFlatDesc 		  = "SELECT id, author, forum, thread, message, parent, isEdited, created FROM posts WHERE thread = $1  ORDER BY id DESC LIMIT $2"
	selectFlatAsc         = "SELECT id, author, forum, thread, message, parent, isEdited, created FROM posts WHERE thread = $1  ORDER BY id ASC LIMIT $2"

	selectParentDesc	  = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created FROM posts as p WHERE p.thread = $1 AND p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 ORDER BY p.path[1] DESC, p.path LIMIT $2)) ORDER BY p.path[1] DESC, p.path"
	selectParentAsc		  = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.isEdited, p.created FROM posts as p WHERE p.thread = $1 AND p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 ORDER BY p.path[1] ASC, p.path LIMIT $2)) ORDER BY p.path[1] ASC, p.path"
	)

func NewPostRepository(db *pgx.ConnPool) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (p *PostRepository) Create(slugOrId string, posts *domain.Posts) error {
	var forum string
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		err = p.db.QueryRow("SELECT id, forum FROM threads WHERE slug = $1 LIMIT 1", slugOrId).Scan(
			&id,
			&forum,
			)
		if err != nil {
			fmt.Println(err, 50)
			return domain.NotFoundError
		}
	} else {
		err = p.db.QueryRow("SELECT forum FROM threads WHERE id = $1 LIMIT 1", id).Scan(
			&forum,
		)
		if err != nil {
			fmt.Println(err, 58)
			return domain.NotFoundError
		}
	}

	if (*posts).Size() == 0 {
		return nil
	}

	if (*posts)[0].Parent != 0 {
		var parent int
		err = p.db.QueryRow("SELECT thread FROM posts WHERE id = $1", (*posts)[0].Parent).Scan(
			&parent,
			)
		if err != nil {
			fmt.Println(err, 73)
			return domain.ParentError
		}
		if parent != id {
			fmt.Println(err, 77)
			return domain.ParentError
		}
	}

	query := "INSERT INTO posts (author, forum, message, parent, thread) values "
	for i, post := range *posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("('%s', '%s', '%s', %d, %d) ", post.Author, forum, post.Message,
			post.Parent, id)
	}

	query += "RETURNING id, created"
	rows, err := p.db.Query(query)
	if err != nil {
		fmt.Println(err, 94)
		return err
	}

	for idx := 0; rows.Next(); idx++ {
		(*posts)[idx].Forum = forum
		(*posts)[idx].Thread = id
		if err = rows.Scan(&(*posts)[idx].Id, &(*posts)[idx].Created); err != nil {
			fmt.Println(err, 102)
			return err
		}
	}

	if rows.Err() != nil {
		pgerr, _ := rows.Err().(pgx.PgError)
		fmt.Println(pgerr, 109)
		switch pgerr.Code {
		case "23505":
			return domain.ConflictError
		default:
			return domain.NotFoundError
		}
	}

	return nil
}

func (p *PostRepository) UpdateById(id int, post *domain.Post) error {
	var query string
	if post.Message == "" {
		query = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE id = $1"

		err := p.db.QueryRow(query, id).Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)

		if err != nil {
			return domain.NotFoundError
		}

		return nil
	}

	query = "UPDATE posts SET message = $1, isEdited = CASE WHEN message = $1 THEN isEdited ELSE true END WHERE id = $2 RETURNING id, parent, author, message, isEdited, forum, thread, created"

	err := p.db.QueryRow(query, post.Message, id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)

	if err != nil {
		return domain.NotFoundError
	}

	return nil

}

func (p *PostRepository) Details(id int, related string) (domain.PostFull, error) {
	postFull := domain.PostFull{}

	post := domain.Post{}

	err := p.db.QueryRow("SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts " +
		"WHERE id = $1", id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
		)

	if err != nil {
		return postFull, domain.NotFoundError
	}

	if strings.Contains(related, "user") {
		user := domain.User{}
		err = p.db.QueryRow("SELECT nickname, fullname, about, email " +
			"FROM users WHERE nickname = $1", post.Author).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)
		if err != nil {
			return postFull, domain.NotFoundError
		}

		postFull.Author = &user
	}

	if strings.Contains(related, "thread") {
		thread := domain.Thread{}

		var readedSlug *string

		err = p.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created " +
			"FROM threads WHERE id = $1 LIMIT 1", post.Thread).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&readedSlug,
			&thread.Created,
		)

		if readedSlug != nil {
			thread.Slug = *readedSlug
		}

		if err != nil {
			return postFull, domain.NotFoundError
		}

		postFull.Thread = &thread
	}

	if strings.Contains(related, "forum") {
		forum := &domain.Forum{}
		err = p.db.QueryRow("SELECT title, author, slug, posts, threads FROM forums " +
			"WHERE slug = $1", post.Forum).Scan(
			&forum.Title,
			&forum.User,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads,
		)

		if err != nil {
			return postFull, domain.NotFoundError
		}

		postFull.Forum = forum
	}

	postFull.Post = &post

	return postFull, nil
}

func (p *PostRepository) GetThreadPosts(slugOrId string, limit int,
	since int, sort string, desc bool) (domain.Posts, error) {
	posts := domain.Posts{}

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		err = p.db.QueryRow("SELECT id FROM threads WHERE slug = $1 LIMIT 1", slugOrId).Scan(
			&id,
			)
		if err != nil {
			return posts, domain.NotFoundError
		}
	} else {
		err = p.db.QueryRow("SELECT id FROM threads WHERE id = $1 LIMIT 1", id).Scan(
			&id,
		)
		if err != nil {
			return posts, domain.NotFoundError
		}
	}


	var query string
	rows := &pgx.Rows{}
	if since == 0 {
		if desc {
			switch sort {
			case "tree" :
				query = selectTreeDesc
			case "parent_tree":
				query = selectParentDesc
			default:
				query = selectFlatDesc
			}
		} else {
			switch sort {
			case "tree" :
				query = selectTreeAsc
			case "parent_tree":
				query = selectParentAsc
			default:
				query = selectFlatAsc
			}
		}
		rows, err = p.db.Query(query, id, limit)
	} else {
		if desc {
			switch sort {
			case "tree" :
				query = selectTreeSinceDesc
			case "parent_tree":
				query = selectParentSinceDesc
			default:
				query = selectFlatSinceDesc
			}
		} else {
			switch sort {
			case "tree" :
				query = selectTreeSinceAsc
			case "parent_tree":
				query = selectParentSinceAsc
			default:
				query = selectFlatSinceAsc
			}
		}
		rows, err = p.db.Query(query, id, since, limit)
	}

	if err != nil {
		return posts, err
	}

	for rows.Next() {
		p := &domain.Post{}
		err := rows.Scan(&p.Id, &p.Author, &p.Forum, &p.Thread, &p.Message, &p.Parent, &p.IsEdited, &p.Created)
		if err != nil {
			return posts, err
		}
		posts = append(posts, *p)
	}
	return posts, nil
}
