package service

import (
	"errors"
	"time"

	"github.com/financeku/backend/internal/config"
	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
	"github.com/financeku/backend/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := hash.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (*models.User, *TokenPair, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, nil, errors.New("account is deactivated")
	}

	if !hash.CheckPassword(input.Password, user.Password) {
		return nil, nil, errors.New("invalid email or password")
	}

	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, errors.New("failed to generate tokens")
	}

	return user, tokens, nil
}

func (s *AuthService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	userID, _ := claims["user_id"].(string)
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) GetCurrentUser(userID string) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) generateTokenPair(user *models.User) (*TokenPair, error) {
	accessExpMinutes := s.cfg.JWTAccessExpMinutes
	refreshExpDays := s.cfg.JWTRefreshExpDays

	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"type":    "access",
		"exp":     time.Now().Add(time.Duration(accessExpMinutes) * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Duration(refreshExpDays) * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessExpMinutes * 60,
	}, nil
}
