package user

import (
	"wstester/internal/entity"
	"wstester/internal/repo/user"
)

type UserService struct {
	userRepo *user.UserRepo
}

func NewUserService(userRepo *user.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) AddUser(user *entity.User) error {
	u.userRepo.
}
