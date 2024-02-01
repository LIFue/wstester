package user

import (
	"wstester/internal/base/data"
	"wstester/internal/entity"
)

type UserRepo struct {
	data *data.Data
}

func NewUserRepo(data *data.Data) *UserRepo {
	return &UserRepo{
		data: data,
	}
}

func (u *UserRepo) Insert(user *entity.User) error {
	return u.data.DB.Create(user).Error
}

func (u *UserRepo) ExistSameUser(user *entity.User) error {
	var userCount int64
	u.data.DB.Model(&entity.User{}).Where(user).Count(&userCount)
}

func (u *UserRepo) FetchPlatformOperator(platform *entity.Platform) {

}
