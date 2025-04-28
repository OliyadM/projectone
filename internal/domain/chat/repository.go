package chat

import "context"

type Repository interface {
	SendMessage(ctx context.Context, msg *ChatMessage) error
	GetMessagesBetweenUsers(ctx context.Context, user1, user2 string) ([]*ChatMessage, error)
	MarkAsSeen(ctx context.Context, messageID string) error
	ListConversationsForUser(ctx context.Context, userID string) ([]*ChatMessage, error)
}
