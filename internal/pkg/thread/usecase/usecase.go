package usecase

import "TPForum/internal/pkg/domain"

type ThreadUsecase struct {
	threadRepo domain.ThreadRepository
}

func NewThreadUsecase(t domain.ThreadRepository) ThreadUsecase {
	return ThreadUsecase{
		threadRepo: t,
	}
}

func (t ThreadUsecase) Create(thread *domain.Thread) error {
	return t.threadRepo.Create(thread)
}

func (t ThreadUsecase) GetThreadBySlug(slug string) (domain.Thread, error) {
	return t.threadRepo.GetThreadBySlug(slug)
}

func (t ThreadUsecase) GetForumThreads(slug string, limit int, since string, desc bool) (domain.Threads, error) {
	return t.threadRepo.GetForumThreads(slug, limit, since, desc)
}

func (t ThreadUsecase) GetThreadBySlugOrId(slugOrId string) (domain.Thread, error) {
	return t.threadRepo.GetThreadBySlugOrId(slugOrId)
}

func (t ThreadUsecase) UpdateThreadBySlugOrId(slugOrId string, thread *domain.Thread) error {
	return t.threadRepo.UpdateThreadBySlugOrId(slugOrId, thread)
}

func (t ThreadUsecase) Vote(slugOrId string, vote domain.Vote) (domain.Thread, error) {
	return t.threadRepo.Vote(slugOrId, vote)
}
