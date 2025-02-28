package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequst struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MessageRequest struct {
	ChatID     uint   `json:"chatid"`
	SenderID   uint   `json:"senderid"`
	ReceiverID uint   `json:"receiverid"`
	Content    string `json:"content"`
}
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
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
	UserID  uint  `json:"userid,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

func (u *User) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw)) == nil
}
