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

func (u *UserService) AddUser(user *entity.User) (err error) {
	var isExist bool
	if isExist, err = u.userRepo.ExistSameUser(user); err != nil {
		return
	}
	u.userRepo.Insert(user)
}
