package dao

import (
	"testlake/model"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type UserDao struct {
	Limit int
}

func NewUserDao() *UserDao {
	return &UserDao{Limit: 50}
}

func (dao *UserDao) Create(user *model.User) error {
	return Database.Create(user).Error
}

func (dao *UserDao) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := Database.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := Database.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := Database.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) GetAll(page int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := Database.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Offset(offset).Limit(dao.Limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (dao *UserDao) Update(user *model.User) error {
	return Database.Save(user).Error
}

func (dao *UserDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.User{}, "id = ?", id).Error
}

func (dao *UserDao) UpdateLastLogin(id uuid.UUID) error {
	return Database.Model(&model.User{}).Where("id = ?", id).Update("last_login_at", "NOW()").Error
}

func (dao *UserDao) UpdateStatus(id uuid.UUID, status model.UserStatus) error {
	return Database.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (dao *UserDao) EmailExists(email string) (bool, error) {
	var count int64
	err := Database.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (dao *UserDao) UsernameExists(username string) (bool, error) {
	var count int64
	err := Database.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
