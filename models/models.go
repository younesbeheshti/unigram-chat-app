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

type ContactsResponse struct {
	Contacts []*User `json:"contacts"`
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
	Messages []*MessageRequest `json:"messages"`
}
type MessageRequest struct {
	ChatID     uint  `json:"chatid,omitempty"`
	SenderName string `json:"sender_name,omitempty"`
	SenderID   uint   `json:"senderid"`
	ReceiverID uint  `json:"receiverid,omitempty"`
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

type Channel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatorID uint      `gorm:"not null" json:"creator_id"`
	Members   []*User   `gorm:"many2many:user_channels;constraint:OnDelete:CASCADE;" json:"members"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Chat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	User1ID   uint      `gorm:"not null;index" json:"user1id"`
	User2ID   uint      `gorm:"not null;index" json:"user2id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID     uint     `gorm:"index" json:"chatid"`
	SenderID   uint      `gorm:"not null" json:"senderid"`
	ReceiverID uint     `gorm:"index" json:"receiverid"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Seen       bool      `gorm:"default:false" json:"seen"`
	CreatedAt  time.Time `json:"created_at"`
}

type Response struct {
	UserID  uint   `json:"userid,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

func (u *User) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw)) == nil
}
