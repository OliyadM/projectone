package user

import "time"

type Role string

const (
	RoleSupplier Role = "supplier"
	RoleReseller Role = "reseller"
	RoleConsumer Role = "consumer"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID              string    `bson:"_id"`
	Name            string    `bson:"name"`
	Username        string    `bson:"username"`
	Email           string    `bson:"email"`
	Password        string    `bson:"password"`
	Role            string    `bson:"role"`
	CreatedAt       time.Time `bson:"created_at"`
	TrustScore      int       `bson:"trust_score"`
	TrustRatedCount int       `bson:"trust_rated_count"` // Total number of rated items
	TrustTotalError float64   `bson:"trust_total_error"`
	IsDeleted       bool      `bson:"is_deleted"`
	IsBlacklisted   bool      `bson:"is_blacklisted"`
	ImageURL        string    `bson:"image_url"` // URL to the user's profile image
}
