package storage

func SaveMessage(chatID uint, senderID uint, receiverID uint, content string) {}
func GetChatHistory(chatID uint)                                              {}
func MarkMessageAsRead(messageID uint)                                        {}
