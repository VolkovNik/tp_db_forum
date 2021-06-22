package repository

import (
	"TPForum/internal/pkg/domain"
	"github.com/jackc/pgx"
)


type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(user domain.User) (*domain.User, error)  {
	_, err := u.db.Exec("INSERT INTO users (nickname, fullname, about, email)" +
		"VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	)

	if err != nil {
		return nil, domain.ConflictError
	}
	return nil, nil
}

func (u *UserRepository) SelectByEmailOrNickname(nickname string, email string) (domain.Users, error) {
	users := domain.Users{}
	rows, err := u.db.Query("SELECT nickname, fullname, about, email FROM users "+
		"WHERE nickname = $1 OR email = $2 LIMIT 2", nickname, email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepository) GetProfileInfo(nickname string) (domain.User, error)  {
	user := domain.User{}
	err := u.db.QueryRow("SELECT nickname, fullname, about, email " +
		"FROM users WHERE nickname = $1", nickname).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
			)

	if err != nil {
		return user, domain.NotFoundError
	}

	return user, nil
}

func (u *UserRepository) UpdateProfileInfo(user *domain.User) error {
	query := "UPDATE users SET   "
	if user.About != "" {
		query += ` about = '` + user.About + "' , "
	}
	if user.Fullname != "" {
		query += " fullname = '" + user.Fullname + "' , "
	}
	if user.Email != "" {
		query += " email = '" + user.Email + "' , "
	}

	query = query[:len(query)-2]

	query += "WHERE nickname = '" + user.Nickname + "' RETURNING about, email, fullname, nickname "
	if user.About == "" && user.Email == "" && user.Fullname == "" {
		query = "SELECT about, email, fullname, nickname FROM users WHERE nickname = '" + user.Nickname + "' "
	}

	err := u.db.QueryRow(query).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return domain.NotFoundError
		default:
			return domain.ConflictError
		}
	}

	return nil
}