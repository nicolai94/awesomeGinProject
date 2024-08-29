package repository

import (
	"awesomeProject/app/domain/dao"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(email string) (dao.User, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func (u AuthRepositoryImpl) FindUserByEmail(email string) (dao.User, error) {
	user := dao.User{
		Email: email,
	}
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Error("Got and error when find user by email. Error: ", err)
		return dao.User{}, err
	}

	return user, nil
}

func AuthRepositoryInit(db *gorm.DB) *AuthRepositoryImpl {
	db.AutoMigrate(&dao.User{})
	return &AuthRepositoryImpl{
		db: db,
	}
}
