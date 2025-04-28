package product

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProduct_GenerateID(t *testing.T) {
	p := &Product{}
	id := p.GenerateID()
	assert.NotEmpty(t, id)
	assert.Len(t, id, 24) // MongoDB ObjectID hex length
}

func TestProduct_ValidateRating(t *testing.T) {
	tests := []struct {
		name    string
		rating  float64
		wantErr bool
	}{
		{"valid rating", 4.5, false},
		{"minimum rating", 0.0, false},
		{"maximum rating", 5.0, false},
		{"negative rating", -1.0, true},
		{"above maximum rating", 5.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{Rating: tt.rating}
			err := p.ValidateRating()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProduct_NewProduct(t *testing.T) {
	resellerID := primitive.NewObjectID()
	supplierID := primitive.NewObjectID().Hex()
	bundleID := primitive.NewObjectID().Hex()

	p := &Product{
		ResellerID:  resellerID,
		SupplierID:  supplierID,
		BundleID:    bundleID,
		Title:       "Test Product",
		Description: "Test Description",
		Size:        "M",
		Type:        "T-Shirt",
		Grade:       "A",
		Price:       29.99,
		Status:      "available",
		ImageURL:    "http://example.com/image.jpg",
		Rating:      4.5,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	// Generate ID before making assertions
	p.ID = p.GenerateID()

	assert.NotEmpty(t, p.ID)
	assert.Equal(t, resellerID, p.ResellerID)
	assert.Equal(t, supplierID, p.SupplierID)
	assert.Equal(t, bundleID, p.BundleID)
	assert.Equal(t, "Test Product", p.Title)
	assert.Equal(t, 4.5, p.Rating)
	assert.Equal(t, "available", p.Status)
}

func TestProduct_UpdateRating(t *testing.T) {
	p := &Product{
		Rating: 4.0,
	}

	// Test valid rating update
	err := p.UpdateRating(4.5)
	assert.NoError(t, err)
	assert.Equal(t, 4.5, p.Rating)

	// Test invalid rating update
	err = p.UpdateRating(5.1)
	assert.Error(t, err)
	assert.Equal(t, 4.5, p.Rating) // Rating should remain unchanged
}
