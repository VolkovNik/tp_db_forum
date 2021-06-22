package usecase

import (
	"TPForum/internal/pkg/domain"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(u domain.UserRepository) UserUsecase {
	return UserUsecase{
		userRepo: u,
	}
}

func (u *UserUsecase) Create(user domain.User) (*domain.User, error)  {
	return u.userRepo.Create(user)
}

func (u *UserUsecase) SelectByEmailOrNickname(nickname string, email string) (domain.Users, error){
	return u.userRepo.SelectByEmailOrNickname(nickname, email)
}

func (u *UserUsecase) GetProfileInfo(nickname string) (domain.User, error) {
	return u.userRepo.GetProfileInfo(nickname)
}

func (u *UserUsecase) UpdateProfileInfo(user *domain.User) error {
	return u.userRepo.UpdateProfileInfo(user)
}