package utils

type Message struct {
	userID int
	content string
}

func NewMsg(userID int, content string) Message {
	return Message{
		userID:userID,
		content: content,
	}
}

func (m *Message) GetUserID() int {
	return m.userID
}