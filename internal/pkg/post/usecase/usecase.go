package usecase

import "TPForum/internal/pkg/domain"

type PostUsecase struct {
	postRepo domain.PostRepository
}

func NewPostUsecase(p domain.PostRepository) PostUsecase {
	return PostUsecase{
		postRepo: p,
	}
}

func (p PostUsecase) Create(slugOrId string, posts *domain.Posts) error  {
	return p.postRepo.Create(slugOrId, posts)
}

func (p PostUsecase) UpdateById(id int, post *domain.Post) error {
	return p.postRepo.UpdateById(id, post)
}

func (p PostUsecase) Details(id int, related string) (domain.PostFull, error) {
	return p.postRepo.Details(id, related)
}

func (p PostUsecase) GetThreadPosts(slugOrId string, limit int,
	since int, sort string, desc bool) (domain.Posts, error) {
	return p.postRepo.GetThreadPosts(slugOrId, limit, since, sort, desc)
}
