package repository

import (
	"awesomeProject/app/domain/dao"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUser(limit, offset int) ([]dao.User, error)
	FindUserById(id string) (dao.User, error)
	Save(user *dao.User) (dao.User, error)
	DeleteUserById(id string) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u UserRepositoryImpl) FindAllUser(limit, offset int) ([]dao.User, error) {
	users := make([]dao.User, 0, limit)

	var err = u.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		log.Error("Got an error finding all couples. Error: ", err)
		return nil, err
	}

	return users, nil
}

func (u UserRepositoryImpl) FindUserById(id string) (dao.User, error) {
	user := dao.User{
		ID: id,
	}
	err := u.db.First(&user).Error
	if err != nil {
		log.Error("Got and error when find user by id. Error: ", err)
		return dao.User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) Save(user *dao.User) (dao.User, error) {
	var err = u.db.Save(user).Error
	if err != nil {
		log.Error("Got an error when save user. Error: ", err)
		return dao.User{}, err
	}
	return *user, nil
}

func (u UserRepositoryImpl) DeleteUserById(id string) error {
	err := u.db.Delete(&dao.User{}, id).Error
	if err != nil {
		log.Error("Got an error when delete user. Error: ", err)
		return err
	}
	return nil
}

func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	db.AutoMigrate(&dao.User{})
	return &UserRepositoryImpl{
		db: db,
	}
}
