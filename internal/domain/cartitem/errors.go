package cartitem

// UnavailableItem represents a cart item that failed validation.
type UnavailableItem struct {
	ListingID string `json:"listingId"`
	Title     string `json:"title"`
}

// CheckoutValidationError is returned when one or more cart items are unavailable.
// It implements the error interface.
type CheckoutValidationError struct {
	Message          string            `json:"message"`
	UnavailableItems []UnavailableItem `json:"unavailableItems"`
}

// Error implements the error interface.
func (e *CheckoutValidationError) Error() string {
	return e.Message
}
