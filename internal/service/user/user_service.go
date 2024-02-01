package user

import (
	"wstester/internal/entity"
	"wstester/internal/repo/user"
	"wstester/pkg/log"
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
	log.Infof("user repo: %+v", u.userRepo)
	var isExist bool
	if isExist, err = u.userRepo.ExistSameUser(user); err != nil || isExist {
		return
	}

	return u.userRepo.Insert(user)
}

func (u *UserService) CreateOrReplaceUser(user *entity.User) (err error) {
	var isExist bool
	if isExist, err = u.userRepo.ExistSameUser(user); err != nil {
		return
	}

	if !isExist {
		return u.userRepo.Insert(user)
	}

	return u.userRepo.UpdateUser(user)
}
