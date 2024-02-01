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

func (u *UserRepo) ExistSameUser(user *entity.User) (bool, error) {
	var userCount int64
	if err := u.data.DB.Model(&entity.User{}).Where(user).Count(&userCount).Error; err != nil {
		return false, err
	}
	return userCount > 0, nil
}

func (u *UserRepo) UpdateUser(user *entity.User) error {
	return u.data.DB.Model(&entity.User{}).Where("uid", user.Uid).Updates(user).Error
}

func (u *UserRepo) FetchPlatformOperator(platform *entity.Platform) {

}
