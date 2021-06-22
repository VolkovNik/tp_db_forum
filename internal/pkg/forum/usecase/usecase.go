package usecase

import "TPForum/internal/pkg/domain"

type ForumUsecase struct {
	forumRepo domain.ForumRepository
}

func NewForumUsecase(f domain.ForumRepository) ForumUsecase {
	return ForumUsecase{
		forumRepo: f,
	}
}

func (f ForumUsecase) Create(forum *domain.Forum) (*domain.Forum, error) {
	return f.forumRepo.Create(forum)
}

func (f ForumUsecase) GetForumDetails(slug string) (*domain.Forum, error) {
	return f.forumRepo.GetForumDetails(slug)
}

func (f ForumUsecase) GetForumUsers(slug string, limit int, since string, desc bool) (domain.Users, error) {
	return f.forumRepo.GetForumUsers(slug, limit, since, desc)
}