package cartitem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	orderusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/google/uuid"
)

// UnavailableItem represents a cart item that failed validation.

type cartItemUsecase struct {
	repo        cartitem.Repository
	productRepo product.Repository // Used to fetch product details
	paymentRepo payment.Repository
	orderUC     orderusecase.OrderUseCase
	orderRepo   order.Repository
}

// NewCartItemUsecase creates a new CartItem usecase instance.
// Note: productRepo is used for product lookup and validation during checkout.
func NewCartItemUsecase(repo cartitem.Repository, productRepo product.Repository, paymentRepo payment.Repository, orderUC orderusecase.OrderUseCase, orderRepo order.Repository) cartitem.Usecase {
	return &cartItemUsecase{
		repo:        repo,
		productRepo: productRepo,
		paymentRepo: paymentRepo,
		orderUC:     orderUC,
		orderRepo:   orderRepo,
	}
}

// AddCartItem adds an item to the user's cart after ensuring no duplicates.

func (u *cartItemUsecase) AddCartItem(ctx context.Context, userID string, listingID string) error {
	// Fetch product details using the provided productID (listingID).
	prod, err := u.productRepo.GetProductByID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if prod == nil {
		return errors.New("product not found")
	}
	// Check that the product is available.
	if prod.Status != "available" {
		return fmt.Errorf("product %s is not available", listingID)
	}

	// Build a new CartItem using the product details.
	cartItem := &cartitem.CartItem{
		ID:        uuid.NewString(),
		UserID:    userID,
		ListingID: prod.ID, // Using product.ID as the listing id.
		Title:     prod.Title,
		Price:     prod.Price,
		ImageURL:  prod.ImageURL,
		Grade:     prod.Grade,
		CreatedAt: time.Now(),
	}

	return u.repo.CreateCartItem(ctx, cartItem)
}

// GetCartItems retrieves all cart items for the given user.
func (u *cartItemUsecase) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	return u.repo.GetCartItems(ctx, userID)
}

// RemoveCartItem deletes a specific item from the user's cart.
func (u *cartItemUsecase) RemoveCartItem(ctx context.Context, userID string, listingID string) error {
	return u.repo.DeleteCartItem(ctx, userID, listingID)
}
func (u *cartItemUsecase) CheckoutCart(ctx context.Context, userID string) (*models.CheckoutResponse, error) {
	items, err := u.repo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var total float64
	var checkoutItems []models.CheckoutItemResponse

	for _, item := range items {
		prod, err := u.productRepo.GetProductByID(ctx, item.ListingID)
		if err != nil || prod == nil {
			return nil, fmt.Errorf("product with ListingID %s not found", item.ListingID)
		}
		if prod.Status != "available" {
			return nil, fmt.Errorf("item %q is no longer available", prod.Title)
		}

		total += prod.Price
		checkoutItems = append(checkoutItems, models.CheckoutItemResponse{
			ListingID: prod.ID,
			Title:     prod.Title,
			Price:     prod.Price,
			SellerID:  prod.ResellerID.Hex(),
			Status:    prod.Status,
		})

		// Create order and payment before marking product as sold
		_, payment, err := u.orderUC.PurchaseProduct(ctx, prod.ID, userID, prod.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to create order for product %s: %w", prod.ID, err)
		}

		// Update payment with reseller ID
		payment.ToUserID = prod.ResellerID.Hex()
		if err := u.paymentRepo.RecordPayment(ctx, payment); err != nil {
			return nil, fmt.Errorf("failed to update payment for product %s: %w", prod.ID, err)
		}

		// Mark product as sold and save to database
		if err := u.productRepo.UpdateProduct(ctx, prod.ID, map[string]interface{}{
			"status": "sold",
		}); err != nil {
			return nil, fmt.Errorf("failed to mark product %s as sold: %w", prod.ID, err)
		}
	}

	platformFee := total * 0.02
	netPayable := total - platformFee

	if err := u.repo.ClearCart(ctx, userID); err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(3 * time.Minute)
		// Simulate delivery status update
	}()

	return &models.CheckoutResponse{
		TotalAmount: total,
		Items:       checkoutItems,
		PlatformFee: platformFee,
		NetPayable:  netPayable,
	}, nil
}
func (u *cartItemUsecase) CheckoutSingleItem(ctx context.Context, userID, listingID string) (*models.CheckoutResponse, error) {
	// Get all cart items for the user
	items, err := u.repo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find the specific item
	var targetItem *cartitem.CartItem
	for _, i := range items {
		if i.ListingID == listingID {
			targetItem = i
			break
		}
	}
	if targetItem == nil {
		return nil, errors.New("item not found in cart")
	}

	// Fetch product info
	prod, err := u.productRepo.GetProductByID(ctx, listingID)
	if err != nil || prod == nil {
		return nil, fmt.Errorf("product with ListingID %s not found", listingID)
	}
	if prod.Status != "available" {
		return nil, fmt.Errorf("item %q is no longer available", prod.Title)
	}

	// Calculate fees
	total := prod.Price
	platformFee := total * 0.02
	netPayable := total - platformFee

	// Create order and payment before marking product as sold
	_, payment, err := u.orderUC.PurchaseProduct(ctx, prod.ID, userID, prod.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to create order for product %s: %w", prod.ID, err)
	}

	// Update payment with reseller ID
	payment.ToUserID = prod.ResellerID.Hex()
	if err := u.paymentRepo.RecordPayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment for product %s: %w", prod.ID, err)
	}

	// Mark product as sold and save to database
	if err := u.productRepo.UpdateProduct(ctx, prod.ID, map[string]interface{}{
		"status": "sold",
	}); err != nil {
		return nil, fmt.Errorf("failed to mark product %s as sold: %w", prod.ID, err)
	}

	// Remove the item from cart
	if err := u.repo.DeleteCartItem(ctx, userID, listingID); err != nil {
		return nil, err
	}

	// Simulate order delivery after 3 minutes
	go func() {
		time.Sleep(3 * time.Minute)
		// TODO: Update order status to "Delivered"
	}()

	return &models.CheckoutResponse{
		TotalAmount: total,
		Items: []models.CheckoutItemResponse{
			{
				ListingID: prod.ID,
				Title:     prod.Title,
				Price:     prod.Price,
				SellerID:  prod.ResellerID.Hex(),
				Status:    "sold", // Update status in response
			},
		},
		PlatformFee: platformFee,
		NetPayable:  netPayable,
	}, nil
}
