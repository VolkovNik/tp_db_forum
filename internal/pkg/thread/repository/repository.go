package repository

import (
	"TPForum/internal/pkg/domain"
	"github.com/jackc/pgx"
	"strconv"
)

var (
	selectThreadsByForumSlugSinceDesc 	= "SELECT id, title, author, forum, message, slug, created, votes FROM threads WHERE forum = $1 AND created <= $2 ORDER BY created DESC LIMIT $3"

	selectThreadsByForumSlugSince 		= "SELECT id, title, author, forum, message, slug, created, votes FROM threads WHERE forum = $1 AND created >= $2 ORDER BY created ASC LIMIT $3"

	selectThreadsByForumSlugDesc 		= "SELECT id, title, author, forum, message, slug, created, votes FROM threads WHERE forum = $1 ORDER BY created DESC LIMIT $2"

	selectThreadsByForumSlug 			= "SELECT id, title, author, forum, message, slug, created, votes FROM threads WHERE forum = $1 ORDER BY created ASC LIMIT $2"
)


type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) *ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (t *ThreadRepository) Create(thread *domain.Thread) error {

	err := t.db.QueryRow("SELECT slug FROM forums WHERE slug = $1 LIMIT 1", thread.Forum).Scan(
		&thread.Forum,
		)

	if err != nil {
		return domain.NotFoundError
	}

	var threadSlug *string
	if thread.Slug != "" {
		threadSlug = &thread.Slug
	}

	err = t.db.QueryRow("INSERT INTO threads (title, author, forum, message, slug, created)" +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		threadSlug,
		thread.Created,
		).Scan(
			&thread.Id,
			)

	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23505":
			return domain.ConflictError
		default:
			return domain.NotFoundError
		}
	}

	return err
}

func (t *ThreadRepository) GetThreadBySlug(slug string) (domain.Thread, error) {
	thread := domain.Thread{}

	err := t.db.QueryRow("SELECT id, title, author, forum, message, slug, created, votes " +
		" FROM threads WHERE slug = $1", slug).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Created,
			&thread.Votes,
			)

	if err != nil {
		return thread, err
	}
	return thread, nil
}

func (t *ThreadRepository) GetForumThreads(slug string, limit int, since string, desc bool) (domain.Threads, error) {
	row, err := t.db.Exec("SELECT 1 FROM forums WHERE slug = $1 LIMIT 1", slug)

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
			query = selectThreadsByForumSlugSinceDesc
		} else {
			query = selectThreadsByForumSlugSince
		}
		rows, err = t.db.Query(query, slug, since, limit)
	} else {
		if desc {
			query = selectThreadsByForumSlugDesc
		} else {
			query = selectThreadsByForumSlug
		}
		rows, err = t.db.Query(query, slug, limit)
	}

	if err != nil {
		return nil, err
	}


	threads := domain.Threads{}

	for rows.Next() {
		thread := domain.Thread{}

		var readedSlug *string

		err = rows.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&readedSlug,
			&thread.Created,
			&thread.Votes,
			)

		if readedSlug != nil {
			thread.Slug = *readedSlug
		}

		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func (t *ThreadRepository) GetThreadBySlugOrId(slugOrId string) (domain.Thread, error) {
	thread := domain.Thread{}

	id, err := strconv.Atoi(slugOrId)
	if err != nil {

		var readedSlug *string

		err = t.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1 LIMIT 1", slugOrId).Scan(
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
			return thread, domain.NotFoundError
		}
	} else {

		var readedSlug *string

		err = t.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1 LIMIT 1", id).Scan(
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
			return thread, domain.NotFoundError
		}
	}

	return thread, nil
}

func (t *ThreadRepository) UpdateThreadBySlugOrId(slugOrId string, thread *domain.Thread) error {
	query := "UPDATE threads SET   "
	if thread.Title != "" {
		query += ` title = '` + thread.Title + "' , "
	}
	if thread.Message != "" {
		query += " message = '" + thread.Message + "' , "
	}

	query = query[:len(query)-2]

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
		query += "WHERE slug = '" + slugOrId + "' RETURNING id, title, author, forum, message, votes, slug, created "
	} else {
		thread.Id = id
		query += "WHERE id = " + slugOrId + " RETURNING id, title, author, forum, message, votes, slug, created "
	}

	if thread.Title == "" && thread.Message == "" {
		query = "SELECT id, title, author, forum, message, votes, slug, created FROM threads "
		if err != nil {
			query += " WHERE slug = '" + slugOrId + "'"
		} else {
			query += " WHERE id = " + slugOrId
		}
	}

	var readedSlug *string

	err = t.db.QueryRow(query).Scan(
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
		return domain.NotFoundError
	}

	return nil
}

func (t *ThreadRepository) Vote(slugOrId string, vote domain.Vote) (domain.Thread, error) {

	thread := domain.Thread{}

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		err = t.db.QueryRow("SELECT id FROM threads WHERE slug = $1 LIMIT 1", slugOrId).Scan(&id)
		if err != nil {
			return thread, domain.NotFoundError
		}
	}

	thread.Id = id

	_, err = t.db.Exec(`INSERT INTO votes (thread_id, nickname, vote)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, nickname) DO UPDATE SET vote = $3`,
		thread.Id,
		vote.Nickname,
		vote.Voice,
	)

	if err != nil {
		return thread, err
	}

	var readedSlug *string

	err = t.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1 LIMIT 1", thread.Id).Scan(
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
		return thread, domain.NotFoundError
	}

	return thread, nil
}
