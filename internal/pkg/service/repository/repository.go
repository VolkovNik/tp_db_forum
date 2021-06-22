package repository

import (
	"TPForum/internal/pkg/domain"
	"github.com/jackc/pgx"
)

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (s *ServiceRepository) Clear() error {

	_, err := s.db.Exec("TRUNCATE forums, users, threads, posts, votes, forums_users CASCADE")

	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceRepository) Status() (domain.Service, error) {
	service := domain.Service{}

	err := s.db.QueryRow("SELECT COUNT(*) from forums").Scan(
		&service.Forum,
		)
	if err != nil {
		return domain.Service{}, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) from users").Scan(
		&service.User,
	)
	if err != nil {
		return domain.Service{}, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) from threads").Scan(
		&service.Thread,
	)
	if err != nil {
		return domain.Service{}, err
	}

	err = s.db.QueryRow("SELECT COUNT(*) from posts").Scan(
		&service.Post,
	)
	if err != nil {
		return domain.Service{}, err
	}

	return service, nil
}