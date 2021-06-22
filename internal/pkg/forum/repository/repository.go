package repository

import (
	"TPForum/internal/pkg/domain"
	"github.com/jackc/pgx"
)

var (
	selectUsersByForumSlugSinceDesc = "SELECT nickname, fullname, about, email FROM forums_users INNER JOIN users u ON u.nickname = forums_users.author WHERE slug = $1 AND nickname < $2 ORDER BY nickname DESC LIMIT $3"

	selectUsersByForumSlugSince 	= "SELECT nickname, fullname, about, email FROM forums_users INNER JOIN users u ON u.nickname = forums_users.author WHERE slug = $1 AND nickname > $2 ORDER BY nickname LIMIT $3"

	selectUsersByForumSlugDesc 		= "SELECT nickname, fullname, about, email FROM forums_users INNER JOIN users u ON u.nickname = forums_users.author WHERE slug = $1 ORDER BY nickname DESC LIMIT $2"

	selectUsersByForumSlug 			= "SELECT nickname, fullname, about, email FROM forums_users INNER JOIN users u ON u.nickname = forums_users.author WHERE slug = $1 ORDER BY nickname LIMIT $2"
)

type ForumRepository struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (f *ForumRepository) Create(forum *domain.Forum) (*domain.Forum, error) {
	err := f.db.QueryRow("SELECT nickname from users where nickname = $1 LIMIT 1", forum.User).Scan(
		&forum.User,
		)

	if err != nil {
		return nil, domain.NotFoundError
	}


	_, err = f.db.Exec("INSERT INTO forums (slug, title, author)" +
								"VALUES ($1, $2, $3)",
		forum.Slug,
		forum.Title,
		forum.User,
		)

	if err != nil {
		pgerr, _ := err.(pgx.PgError)
		switch pgerr.Code {
		case "23505":
			return nil, domain.ConflictError
		default:
			return nil, domain.NotFoundError
		}
	} else {
		return nil, nil
	}
}

func (f *ForumRepository) GetForumDetails(slug string) (*domain.Forum, error) {
	forum := &domain.Forum{}
	err := f.db.QueryRow("SELECT title, author, slug, posts, threads FROM forums " +
		"WHERE slug = $1",
		slug).Scan(
				&forum.Title,
				&forum.User,
				&forum.Slug,
				&forum.Posts,
				&forum.Threads,
			)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (f *ForumRepository) GetForumUsers(slug string, limit int, since string, desc bool) (domain.Users, error)  {
	row, err := f.db.Exec("SELECT 1 FROM forums WHERE slug = $1 LIMIT 1", slug)

	if err != nil {
		return nil, err
	}

	if row.RowsAffected() == 0 {
		return nil, domain.NotFoundError
	}

	var query string
	rows := &pgx.Rows{}

	if since != "" {
		if desc {
			query = selectUsersByForumSlugSinceDesc
		} else {
			query = selectUsersByForumSlugSince
		}
		rows, err = f.db.Query(query, slug, since, limit)
	} else {
		if desc {
			query = selectUsersByForumSlugDesc
		} else {
			query = selectUsersByForumSlug
		}
		rows, err = f.db.Query(query, slug, limit)
	}

	if err != nil {
		return nil, err
	}

	users := domain.Users{}

	for rows.Next() {
		user := domain.User{}


		err = rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
