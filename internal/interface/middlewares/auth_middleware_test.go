package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJWTService is a mock implementation of auth.JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID, username, role string) (string, error) {
	args := m.Called(userID, username, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ParseToken(token string) (*jwt.Token, jwt.MapClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*jwt.Token), args.Get(1).(jwt.MapClaims), args.Error(2)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		tokenClaims    jwt.MapClaims
		parseTokenErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Missing or invalid token"}`,
		},
		{
			name:           "Invalid Authorization Format",
			authHeader:     "InvalidFormat token123",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Missing or invalid token"}`,
		},
		{
			name:           "Invalid Token",
			authHeader:     "Bearer invalid-token",
			parseTokenErr:  jwt.ErrSignatureInvalid,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Unauthorized"}`,
		},
		{
			name:           "Valid Token with Invalid UserID",
			authHeader:     "Bearer valid-token",
			tokenClaims:    jwt.MapClaims{"user_id": "invalid-id"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid user ID format in token"}`,
		},
		{
			name:           "Valid Token with Missing UserID",
			authHeader:     "Bearer valid-token",
			tokenClaims:    jwt.MapClaims{"username": "testuser"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or missing user ID in token"}`,
		},
		{
			name:           "Valid Token with Valid Claims",
			authHeader:     "Bearer valid-token",
			tokenClaims: jwt.MapClaims{
				"user_id":  "507f1f77bcf86cd799439011",
				"username": "testuser",
				"role":     "user",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockJWTService := new(MockJWTService)
			mockToken := &jwt.Token{Claims: tt.tokenClaims}
			mockJWTService.On("ParseToken", mock.Anything).Return(mockToken, tt.tokenClaims, tt.parseTokenErr)

			router := setupRouter()
			router.Use(AuthMiddleware(mockJWTService))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestAuthorizeRoles(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		allowedRoles   []string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Authorized Role",
			role:           "admin",
			allowedRoles:   []string{"admin", "user"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "Unauthorized Role",
			role:           "guest",
			allowedRoles:   []string{"admin", "user"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"Access denied: insufficient role"}`,
		},
		{
			name:           "Missing Role",
			role:           "",
			allowedRoles:   []string{"admin", "user"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"Role not found in token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := setupRouter()
			router.Use(func(c *gin.Context) {
				if tt.role != "" {
					c.Set("role", tt.role)
				}
				c.Next()
			})
			router.Use(AuthorizeRoles(tt.allowedRoles...))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
} 