package usecase

import "TPForum/internal/pkg/domain"

type ServiceUsecase struct {
	serviceRepo domain.ServiceRepository
}

func NewServiceUsecase(s domain.ServiceRepository) ServiceUsecase {
	return ServiceUsecase{
		serviceRepo: s,
	}
}

func (s ServiceUsecase) Clear() error {
	return s.serviceRepo.Clear()
}

func (s ServiceUsecase) Status() (domain.Service, error) {
	return s.serviceRepo.Status()
}
