package payment

type PaymentType string

const (
	B2B PaymentType = "b2b"
	B2C PaymentType = "b2c"
)

type Payment struct {
	ID            string
	FromUserID    string
	ToUserID      string
	Amount        float64
	PlatformFee   float64
	SellerEarning float64
	Status        string
	ReferenceID   string      // This is either BundleID or ProductID
	Type          PaymentType // "b2b" or "b2c"
	CreatedAt     string
}
