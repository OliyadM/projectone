package chat

type MessageType string

const (
	Text  MessageType = "text"
	Image MessageType = "image"
	File  MessageType = "file"
)

type ChatMessage struct {
	ID         string
	SenderID   string
	ReceiverID string
	Text       string
	Type       MessageType
	Timestamp  string
	Seen       bool
}
