package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequst struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ContactsRespose struct {
	Contacts *[]User `json:"contacts"`
}

type ChatResponse struct {
	Chats []*Chat `json:"chat"`
}

type AddChatRequest struct {
	User1ID uint `json:"user1id"`
	User2ID uint `json:"user2id"`
}

type AddChatResponse struct {
	ChatID uint `json:"chat_id"`
}

type MessageHistory struct {
	Messages []*Message `json:"messages"`
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
	ID        uint   `gorm:"primarykey" json:"id"`
	User1ID   uint   `gorm:"not null" json:"user1id"`
	User2ID   uint   `gorm:"not null" json:"user2id"`
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
	UserID  uint   `json:"userid,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

func (u *User) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw)) == nil
}
