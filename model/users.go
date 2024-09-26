package model

import (
	"encoding/json"
	"fmt"
	"go-fiber/database"
	"log"
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
type LoginUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UsersRedis struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
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

func (user *Users) AfterUpdate(tx *gorm.DB) (err error) {
	// Convert user data to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshalling user data:", err)
		return
	}

	// Set the updated user in Redis
	if err := database.SetKey(tx.Statement.Context, "users:"+fmt.Sprint(user.ID), userJSON); err != nil {
		log.Println("Error setting user in Redis:", err)
		return err
	}

	log.Println("User updated and synced to Redis:", user.ID)
	return nil
}
