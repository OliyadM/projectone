package reviewusecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"github.com/google/uuid"
)

type reviewUsecase struct {
	reviewRepo review.Repository
	orderRepo  order.Repository
}

func NewReviewUsecase(reviewRepo review.Repository, orderRepo order.Repository) review.Usecase {
	return &reviewUsecase{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
	}
}

func (u *reviewUsecase) SubmitReview(ctx context.Context, r *review.Review) error {
	fmt.Printf("🔍 Starting review submission process\n")
	fmt.Printf("📝 Review details: %+v\n", r)

	// Check if the order exists and is delivered
	fmt.Printf("🔍 Checking order status for ID: %s\n", r.OrderID)
	order, err := u.orderRepo.GetOrderByID(ctx, r.OrderID)
	if err != nil {
		fmt.Printf("❌ Error fetching order: %v\n", err)
		return errors.New("order not found: " + err.Error())
	}
	if order == nil {
		fmt.Printf("❌ Order not found\n")
		return errors.New("order not found")
	}
	fmt.Printf("✅ Found order: %+v\n", order)

	if order.Status != "completed" {
		fmt.Printf("❌ Order not delivered. Current status: %s\n", order.Status)
		return errors.New("cannot review before delivery")
	}

	// Check if the user already reviewed this product
	fmt.Printf("🔍 Checking for existing review by user %s for product %s\n", r.UserID, r.ProductID)
	existingReview, err := u.reviewRepo.GetReviewByUserAndProduct(ctx, r.UserID, r.ProductID)
	if err != nil {
		fmt.Printf("❌ Error checking existing review: %v\n", err)
		return err
	}
	if existingReview != nil {
		fmt.Printf("❌ User already reviewed this product\n")
		return errors.New("you already reviewed this item")
	}

	// Save the review
	r.ID = uuid.NewString()
	r.CreatedAt = time.Now()
	fmt.Printf("📝 Saving review with ID: %s\n", r.ID)
	
	if err := u.reviewRepo.CreateReview(ctx, r); err != nil {
		fmt.Printf("❌ Error saving review: %v\n", err)
		return err
	}

	fmt.Printf("✅ Review saved successfully\n")
	return nil
}

func (u *reviewUsecase) GetResellerReviews(ctx context.Context, resellerID string) ([]*review.Review, error) {
	fmt.Printf("🔍 Fetching reviews for reseller: %s\n", resellerID)
	
	reviews, err := u.reviewRepo.GetReviewsByReseller(ctx, resellerID)
	if err != nil {
		fmt.Printf("❌ Error fetching reseller reviews: %v\n", err)
		return nil, fmt.Errorf("failed to fetch reseller reviews: %w", err)
	}

	fmt.Printf("✅ Found %d reviews for reseller\n", len(reviews))
	return reviews, nil
}
