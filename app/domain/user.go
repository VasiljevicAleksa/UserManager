package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	UniqueConstraintNickname = "users_nickname_key"
	UniqueConstraintEmail    = "users_email_key"
)

type User struct {
	Id        uuid.UUID `gorm:"column:id;primaryKey"`
	Firstname string    `gorm:"column:first_name;not null"`
	Lastname  string    `gorm:"column:last_name;not null"`
	Nickname  string    `gorm:"column:nickname;unique;not null"`
	Password  string    `gorm:"column:password;not null"`
	Email     string    `gorm:"column:email;unique;not null"`
	Country   string    `gorm:"column:country;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
