package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/kaori/backend/internal/config"
	"github.com/kaori/backend/internal/dummy"
	"github.com/kaori/backend/internal/model"
	jwtutil "github.com/kaori/backend/pkg/jwt"
)

// AuthService handles authentication logic
type AuthService struct {
	cfg *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

// Dummy passwords (in production these would be hashed)
var dummyPasswords = map[string]string{
	"admin@kaori.pos":   "admin123",
	"store@kaori.pos":   "store123",
	"cashier@kaori.pos": "cashier123",
	"kitchen@kaori.pos": "kitchen123",
}

// Login authenticates a user with email and password
func (s *AuthService) Login(email, password string) (*model.LoginResponse, error) {
	// Find user in dummy data
	var dummyUser *dummy.User
	for i := range dummy.Users {
		if dummy.Users[i].Email == email {
			dummyUser = &dummy.Users[i]
			break
		}
	}

	if dummyUser == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	expectedPassword, ok := dummyPasswords[email]
	if !ok || password != expectedPassword {
		return nil, errors.New("invalid email or password")
	}

	return s.generateTokens(dummyUser)
}

// LoginWithPIN authenticates a user with email and PIN
func (s *AuthService) LoginWithPIN(email, pin string) (*model.LoginResponse, error) {
	var dummyUser *dummy.User
	for i := range dummy.Users {
		if dummy.Users[i].Email == email && dummy.Users[i].PIN == pin {
			dummyUser = &dummy.Users[i]
			break
		}
	}

	if dummyUser == nil {
		return nil, errors.New("invalid email or PIN")
	}

	return s.generateTokens(dummyUser)
}

// RefreshToken refreshes an access token (simplified for dummy data)
func (s *AuthService) RefreshToken(refreshToken string) (*model.LoginResponse, error) {
	// For dummy data, just return the first admin user
	for i := range dummy.Users {
		if dummy.Users[i].Role == "super_admin" {
			return s.generateTokens(&dummy.Users[i])
		}
	}
	return nil, errors.New("invalid refresh token")
}

// Logout invalidates refresh tokens for a user
func (s *AuthService) Logout(userID string) error {
	// No-op for dummy data
	return nil
}

// GetCurrentUser returns the current user by ID
func (s *AuthService) GetCurrentUser(userID string) (*model.User, error) {
	for _, u := range dummy.Users {
		if u.ID == userID {
			return &model.User{
				ID:    uuid.MustParse(u.ID),
				Email: u.Email,
				Name:  u.Name,
				Role:  u.Role,
			}, nil
		}
	}
	return nil, errors.New("user not found")
}

// generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(dummyUser *dummy.User) (*model.LoginResponse, error) {
	// Generate access token
	accessToken, err := jwtutil.GenerateToken(
		dummyUser.ID,
		dummyUser.Email,
		dummyUser.Role,
		"store-1", // Default store
		s.cfg.JWTSecret,
		s.cfg.JWTExpiryHours,
	)
	if err != nil {
		return nil, err
	}

	// Generate simple refresh token
	refreshTokenStr := jwtutil.GenerateRefreshToken()

	// Convert to model.User
	user := model.User{
		ID:    uuid.MustParse(dummyUser.ID),
		Email: dummyUser.Email,
		Name:  dummyUser.Name,
		Role:  dummyUser.Role,
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    s.cfg.JWTExpiryHours * 3600,
		User:         user,
	}, nil
}
