package usecase

import (
	"context"
	"travel_advisor/domain"
)

type UserUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecase{
		userRepository: userRepo,
	}
}

func (u *UserUsecase) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	return u.userRepository.Create(ctx, user)
}

func (u *UserUsecase) Get(ctx context.Context, ctr *domain.UserCriteria) (*domain.User, error) {
	return u.userRepository.Get(ctx, ctr)
}
