package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Register(ctx context.Context, user user.User) (*auth.LoginResult, error) {
	args := m.Called(ctx, user)
	result, _ := args.Get(0).(*auth.LoginResult)
	return result, args.Error(1)
}

func (m *MockAuthUsecase) Login(ctx context.Context, creds auth.LoginCredentials) (*auth.LoginResult, error) {
	args := m.Called(ctx, creds)

	// Return a typed pointer for LoginResult
	result, _ := args.Get(0).(*auth.LoginResult)
	return result, args.Error(1)
}

type AuthControllerTestSuite struct {
	suite.Suite
	controller *AuthController
	mockUC     *MockAuthUsecase
	router     *gin.Engine
}

func (suite *AuthControllerTestSuite) SetupTest() {
	suite.mockUC = new(MockAuthUsecase)
	suite.controller = NewAuthController(suite.mockUC)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func (suite *AuthControllerTestSuite) TestRegister_Success() {
	// Setup
	newUser := user.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     string(user.RoleConsumer),
	}
	loginResult := &auth.LoginResult{
		Token:    "test-token",
		ID:       "user-id-123",
		Username: "test@example.com",
		Role:     "consumer",
	}

	suite.mockUC.On("Register", mock.Anything, newUser).Return(loginResult, nil)

	// Execute
	jsonData, _ := json.Marshal(newUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	suite.router.POST("/auth/register", suite.controller.Register)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), loginResult.Token, response["token"])
	user := response["user"].(map[string]interface{})
	assert.Equal(suite.T(), loginResult.ID, user["id"])
	assert.Equal(suite.T(), loginResult.Username, user["username"])
	assert.Equal(suite.T(), loginResult.Role, user["role"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AuthControllerTestSuite) TestRegister_InvalidRequest() {
	// Setup
	invalidUser := "invalid json"

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(invalidUser)))
	req.Header.Set("Content-Type", "application/json")
	suite.router.POST("/auth/register", suite.controller.Register)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AuthControllerTestSuite) TestRegister_Conflict() {
	// Setup
	newUser := user.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     string(user.RoleConsumer),
	}
	errorMsg := "user already exists"

	suite.mockUC.On("Register", mock.Anything, newUser).Return("", errors.New(errorMsg))

	// Execute
	jsonData, _ := json.Marshal(newUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	suite.router.POST("/auth/register", suite.controller.Register)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), errorMsg, response["error"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AuthControllerTestSuite) TestLogin_Success() {
	// Setup
	creds := auth.LoginCredentials{
		Username: "test@example.com",
		Password: "password123",
	}
	loginResult := &auth.LoginResult{
		Token:    "test-token",
		ID:       "user-id-123",
		Username: "test@example.com",
		Role:     "consumer",
	}
	suite.mockUC.On("Login", mock.Anything, creds).Return(loginResult, nil)

	// Execute
	jsonData, _ := json.Marshal(creds)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	suite.router.POST("/auth/login", suite.controller.Login)
	suite.router.ServeHTTP(w, req)

	// Assert
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(suite.T(), loginResult.Token, response["token"])
	user := response["user"].(map[string]interface{})
	assert.Equal(suite.T(), loginResult.ID, user["id"])
	assert.Equal(suite.T(), loginResult.Username, user["username"])
	assert.Equal(suite.T(), loginResult.Role, user["role"])
}

func (suite *AuthControllerTestSuite) TestLogin_InvalidRequest() {
	// Setup
	invalidCreds := "invalid json"

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(invalidCreds)))
	req.Header.Set("Content-Type", "application/json")
	suite.router.POST("/auth/login", suite.controller.Login)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AuthControllerTestSuite) TestLogin_Unauthorized() {
	// Setup
	creds := auth.LoginCredentials{
		Username: "test@example.com",
		Password: "wrongpassword",
	}
	errorMsg := "invalid username or password"

	suite.mockUC.On("Login", mock.Anything, creds).Return((*auth.LoginResult)(nil), errors.New(errorMsg))

	// Execute
	jsonData, _ := json.Marshal(creds)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	suite.router.POST("/auth/login", suite.controller.Login)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), errorMsg, response["error"])
	suite.mockUC.AssertExpectations(suite.T())
}

func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
