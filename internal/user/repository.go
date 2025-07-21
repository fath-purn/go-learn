package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]User, error)
	FindByID(ID int) (User, error)
	FindByEmail(email string) (User, error)
	RegisterUser(user User) (User, error)
	LoginUser(user User) (User, error)
	Update(user User) (User, error)
	Delete(ID int) error
	FindByVerificationToken(token string) (User, error)
	ResetPassword(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) FindByID(ID int) (User, error) {
	var user User
	if err := r.db.First(&user, ID).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	if err := r.db.Where("email = ? ", email).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) RegisterUser(user User) (User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) LoginUser(user User) (User, error) {
	var foundUser User
	if err := r.db.Where("email = ? ", user.Email).First(&foundUser).Error; err != nil {
		return User{}, err
	}
	return foundUser, nil
}

func (r *repository) Update(user User) (User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) Delete(ID int) error {
	if err := r.db.Delete(&User{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByVerificationToken(token string) (User, error) {
	var user User
	if err := r.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) ResetPassword(user User) (User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return User{},
			err
	}
	return user, nil
}
