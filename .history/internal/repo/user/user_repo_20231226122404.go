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

func (u *UserRepo) FetchPlatformOperator(platform *entity.Platform) {

}
