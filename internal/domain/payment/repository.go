package payment

import "context"

type Repository interface {
	RecordPayment(ctx context.Context, p *Payment) error
	GetPaymentsByUser(ctx context.Context, userID string) ([]*Payment, error)
	GetPaymentsByType(ctx context.Context, userID string, pType PaymentType) ([]*Payment, error)
	GetAllPlatformFees(ctx context.Context) (float64, float64, error)
}
