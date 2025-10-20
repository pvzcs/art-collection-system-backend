package service

import (
	"art-collection-system/internal/models"
	"art-collection-system/internal/repository"
	"art-collection-system/internal/utils"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo *repository.UserRepository
	redis    *redis.Client
}

// NewAuthService creates a new authentication service instance
func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		redis:    redisClient,
	}
}

// SendVerificationCode generates a 6-digit verification code and stores it in Redis
// Key format: verify:email:{email}, TTL: 5 minutes
func (s *AuthService) SendVerificationCode(email string) error {
	// Generate 6-digit random verification code
	code := generateVerificationCode()

	// Store in Redis with 5 minutes TTL
	ctx := context.Background()
	key := fmt.Sprintf("verify:email:%s", email)
	err := s.redis.Set(ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	// TODO: Send email with verification code
	// For now, we'll just log it (in production, integrate with email service)
	fmt.Printf("Verification code for %s: %s\n", email, code)

	return nil
}

// Register registers a new user after validating email uniqueness, verification code, and encrypting password
// Default role is "user"
func (s *AuthService) Register(email, code, password, nickname string) (*models.User, error) {
	// Validate email uniqueness
	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Validate verification code
	ctx := context.Background()
	key := fmt.Sprintf("verify:email:%s", email)
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("verification code expired or not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve verification code: %w", err)
	}
	if storedCode != code {
		return nil, errors.New("invalid verification code")
	}

	// Encrypt password using Bcrypt
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user record with default role "user"
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Nickname: nickname,
		Role:     "user",
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Delete verification code after successful registration
	s.redis.Del(ctx, key)

	return user, nil
}

// Login validates email and password, then generates a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Validate password
	err = utils.ComparePassword(user.Password, password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// Logout adds the JWT token to Redis blacklist
// Key format: blacklist:{token}, TTL: remaining token validity period
func (s *AuthService) Logout(token string) error {
	// Validate token to get expiration time
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Calculate remaining TTL
	expiresAt := claims.ExpiresAt.Time
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Add token to blacklist
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", token)
	err = s.redis.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// ValidateToken validates JWT token and checks if it's in the blacklist
func (s *AuthService) ValidateToken(token string) (*models.User, error) {
	// Check if token is in blacklist
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", token)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if exists > 0 {
		return nil, errors.New("token has been revoked")
	}

	// Validate token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return user, nil
}

// generateVerificationCode generates a 6-digit random verification code
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // Generate number between 100000 and 999999
	return fmt.Sprintf("%06d", code)
}
