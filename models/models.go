package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primarykey;autoIncreament"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Chat struct {
	ID        uint `gorm:"primarykey;autoIncreament"`
	User1ID   uint `gorm:"not null"`
	User2ID   uint `gorm:"not null"`
	CreatedAt time.Time
}

type Message struct {
	ID         uint   `gorm:"primarykey;autoIncreament"`
	ChatID     uint   `gorm:"not null;index"`
	SenderID   uint   `gorm:"not null"`
	ReceiverID uint   `gorm:"not null"`
	Content    string `gorm:"type:text;not null"`
	Seen       bool   `gorm:"defaul:false"`
	CreatedAt  time.Time
}

type Respnse struct {
	UserID  int64  `json:"userid,omitempty"`
	Message string `json:"message,omitempty"`
}
