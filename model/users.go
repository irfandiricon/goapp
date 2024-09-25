package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Users represents the user model in the database.
type Users struct {
	gorm.Model
	ID        int    `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UsersRedis struct {
	gorm.Model
	ID    int    `gorm:"primaryKey"`
	Email string `gorm:"uniqueIndex"`
	Name  string
}

// HashPassword hashes the user's password.
func (user *Users) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password.
func (user *Users) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
