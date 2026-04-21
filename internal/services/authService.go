package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/repository"
)

type AuthServiceConfig struct {
	JWTSecret           string
	AccessTokenTTL      time.Duration
	MaxLoginAttempts    int
	LoginWindowDuration time.Duration
}

type loginAttempt struct {
	Count       int
	WindowStart time.Time
}

type AuthService struct {
	UserRepo repository.UserRepository
	Config   AuthServiceConfig

	attemptMu sync.Mutex
	attempts  map[string]loginAttempt

	denylistMu sync.RWMutex
	denylist   map[string]time.Time
}



func (s *AuthService) Register(user entity.User) error {
	user.Role = entity.Uploader
	return s.createUser(user)
}



func (s *AuthService) RegisterAdmin(user entity.User) error {
	user.Role = entity.Admin
	return s.createUser(user)
}



func (s *AuthService) Login(username, password string, sourceKey string) (string, error) {
	normalizedUsername := strings.TrimSpace(strings.ToLower(username))
	if normalizedUsername == "" {
		return "", errors.New("username is required")
	}

	if s.isRateLimited(sourceKey) {
		return "", errors.New("too many failed login attempts, retry later")
	}

	user, err := s.UserRepo.GetUserByUsername(context.Background(), normalizedUsername)
	if err != nil {
		s.recordFailedAttempt(sourceKey)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.recordFailedAttempt(sourceKey)
		return "", errors.New("invalid credentials")
	}

	s.resetFailedAttempt(sourceKey)
	return s.generateToken(user)
}



func (s *AuthService) Logout(token string) error {
	claims, err := s.ValidateToken(token)
	if err != nil {
		return err
	}

	s.denylistMu.Lock()
	defer s.denylistMu.Unlock()
	s.denylist[token] = claims.ExpiresAt.Time
	return nil
}



func NewAuthService(repo repository.UserRepository, cfg AuthServiceConfig) *AuthService {
	if cfg.MaxLoginAttempts <= 0 {
		cfg.MaxLoginAttempts = 5
	}
	if cfg.LoginWindowDuration <= 0 {
		cfg.LoginWindowDuration = 15 * time.Minute
	}
	if cfg.AccessTokenTTL <= 0 {
		cfg.AccessTokenTTL = time.Hour
	}

	return &AuthService{
		UserRepo: repo,
		Config:   cfg,
		attempts: make(map[string]loginAttempt),
		denylist: make(map[string]time.Time),
	}
}



func (s *AuthService) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	s.denylistMu.RLock()
	expiry, blocked := s.denylist[tokenString]
	s.denylistMu.RUnlock()
	if blocked && expiry.After(time.Now().UTC()) {
		return nil, errors.New("token is logged out")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}



func (s *AuthService) ParseRole(tokenString string) (entity.UserRole, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.Config.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token payload")
	}

	rawRole, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("token role is missing")
	}

	return entity.UserRole(rawRole), nil
}



func (s *AuthService) createUser(user entity.User) error {
	normalizedUsername := strings.TrimSpace(strings.ToLower(user.Username))
	if normalizedUsername == "" {
		return errors.New("username is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Username = normalizedUsername
	user.PasswordHash = string(hash)
	return s.UserRepo.CreateUser(context.Background(), user)
}


func (s *AuthService) generateToken(user entity.User) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":      user.Username,
		"username": user.Username,
		"role":     string(user.Role),
		"iat":      now.Unix(),
		"exp":      now.Add(s.Config.AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Config.JWTSecret))
}



func (s *AuthService) isRateLimited(sourceKey string) bool {
	s.attemptMu.Lock()
	defer s.attemptMu.Unlock()

	now := time.Now().UTC()
	attempt, ok := s.attempts[sourceKey]
	if !ok {
		return false
	}

	if now.Sub(attempt.WindowStart) > s.Config.LoginWindowDuration {
		delete(s.attempts, sourceKey)
		return false
	}

	return attempt.Count >= s.Config.MaxLoginAttempts
}



func (s *AuthService) recordFailedAttempt(sourceKey string) {
	s.attemptMu.Lock()
	defer s.attemptMu.Unlock()

	now := time.Now().UTC()
	attempt, ok := s.attempts[sourceKey]
	if !ok || now.Sub(attempt.WindowStart) > s.Config.LoginWindowDuration {
		s.attempts[sourceKey] = loginAttempt{Count: 1, WindowStart: now}
		return
	}

	attempt.Count++
	s.attempts[sourceKey] = attempt
}



func (s *AuthService) resetFailedAttempt(sourceKey string) {
	s.attemptMu.Lock()
	defer s.attemptMu.Unlock()
	delete(s.attempts, sourceKey)
}
