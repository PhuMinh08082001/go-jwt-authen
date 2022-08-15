package repository

import (
	"github.com/PhuMinh08082001/go-jwt-authen/internal/dal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (u *UserRepository) GetUser(userName string) *model.User {
	var user = &model.User{}
	u.Db.Where("user_name = ?", userName).First(&user)

	return user
}
