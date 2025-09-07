package user

import (
	"context"
	"gogym-api/internal/domain/common"
	"gogym-api/internal/domain/user"
	"log/slog"
	"time"
)

// Repository interface for user data access
type Repository interface {
	FindByID(ctx context.Context, id common.ID) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id common.ID) error
}

// RefreshTokenRepository interface for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *user.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*user.RefreshToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteExpiredTokens(ctx context.Context) error
	DeleteAllByUserID(ctx context.Context, userID common.ID) error
}

// PasswordService interface for password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}

// TokenService interface for JWT operations
type TokenService interface {
	GenerateTokens(userID common.ID, email string) (accessToken, refreshToken string, err error)
	ValidateAccessToken(token string) (userID common.ID, email string, err error)
	ValidateRefreshToken(token string) (userID common.ID, tokenHash string, err error)
}

// UseCase represents user use cases
type UseCase struct {
	userRepo         Repository
	refreshTokenRepo RefreshTokenRepository
	passwordService  PasswordService
	tokenService     TokenService
	logger           *slog.Logger
}

// NewUseCase creates a new user use case
func NewUseCase(
	userRepo Repository,
	refreshTokenRepo RefreshTokenRepository,
	passwordService PasswordService,
	tokenService TokenService,
	logger *slog.Logger,
) *UseCase {
	return &UseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordService:  passwordService,
		tokenService:     tokenService,
		logger:           logger,
	}
}

// SignupRequest represents signup input
type SignupRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=128"`
	DisplayName string `json:"display_name" validate:"required,max=100"`
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication output
type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int        `json:"expires_in"` // seconds
	User         *user.User `json:"user"`
}

// Signup creates a new user account
func (uc *UseCase) Signup(ctx context.Context, req SignupRequest) (*AuthResponse, error) {
	uc.logger.InfoContext(ctx, "user signup attempt", "email", req.Email, "display_name", req.DisplayName)

	// Validate password strength
	if err := user.ValidatePassword(req.Password); err != nil {
		uc.logger.WarnContext(ctx, "weak password", "email", req.Email, "error", err)
		return nil, err
	}

	// Check if user already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		uc.logger.WarnContext(ctx, "signup attempt with existing email", "email", req.Email)
		return nil, common.NewDomainError(common.ErrAlreadyExists, "email_exists", "email already registered")
	}

	// Create user entity
	email, err := user.NewEmail(req.Email)
	if err != nil {
		uc.logger.ErrorContext(ctx, "invalid email", "error", err)
		return nil, common.NewDomainError(common.ErrInvalidEmail, "invalid_email", "invalid email format")
	}
	
	newUser, err := user.NewUser(email, req.DisplayName)
	if err != nil {
		uc.logger.ErrorContext(ctx, "invalid user data", "error", err)
		return nil, err
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to hash password", "error", err)
		return nil, common.NewDomainError(common.ErrInternal, "password_hash_failed", "failed to secure password")
	}
	newUser.PasswordHash = hashedPassword

	// Save user
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		uc.logger.ErrorContext(ctx, "failed to create user", "email", req.Email, "error", err)
		return nil, common.NewDomainErrorWithCause(err, "user_creation_failed", "failed to create user")
	}

	// Generate tokens
	accessToken, refreshToken, err := uc.tokenService.GenerateTokens(newUser.ID, newUser.Email.String())
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to generate tokens", "user_id", newUser.ID, "error", err)
		return nil, common.NewDomainError(common.ErrInternal, "token_generation_failed", "failed to generate authentication tokens")
	}

	// Save refresh token
	if err := uc.saveRefreshToken(ctx, newUser.ID, refreshToken); err != nil {
		uc.logger.ErrorContext(ctx, "failed to save refresh token", "user_id", newUser.ID, "error", err)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "user signup successful", "user_id", newUser.ID, "email", req.Email)

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         newUser,
	}, nil
}

// Login authenticates a user
func (uc *UseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	uc.logger.InfoContext(ctx, "user login attempt", "email", req.Email)

	// Find user by email
	foundUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || foundUser == nil {
		uc.logger.WarnContext(ctx, "login attempt with non-existent email", "email", req.Email)
		return nil, common.NewDomainError(common.ErrUnauthorized, "invalid_credentials", "invalid email or password")
	}

	// Verify password
	if err := uc.passwordService.VerifyPassword(req.Password, foundUser.PasswordHash); err != nil {
		uc.logger.WarnContext(ctx, "login attempt with invalid password", "email", req.Email)
		return nil, common.NewDomainError(common.ErrUnauthorized, "invalid_credentials", "invalid email or password")
	}

	// Generate new tokens
	accessToken, refreshToken, err := uc.tokenService.GenerateTokens(foundUser.ID, foundUser.Email.String())
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to generate tokens", "user_id", foundUser.ID, "error", err)
		return nil, common.NewDomainError(common.ErrInternal, "token_generation_failed", "failed to generate authentication tokens")
	}

	// Revoke existing refresh tokens and save new one
	if err := uc.refreshTokenRepo.DeleteAllByUserID(ctx, foundUser.ID); err != nil {
		uc.logger.WarnContext(ctx, "failed to revoke existing tokens", "user_id", foundUser.ID, "error", err)
	}

	if err := uc.saveRefreshToken(ctx, foundUser.ID, refreshToken); err != nil {
		uc.logger.ErrorContext(ctx, "failed to save refresh token", "user_id", foundUser.ID, "error", err)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "user login successful", "user_id", foundUser.ID, "email", req.Email)

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         foundUser,
	}, nil
}

// RefreshToken refreshes access token using refresh token
func (uc *UseCase) RefreshToken(ctx context.Context, refreshTokenStr string) (*AuthResponse, error) {
	uc.logger.InfoContext(ctx, "token refresh attempt")

	// Validate refresh token
	userID, tokenHash, err := uc.tokenService.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		uc.logger.WarnContext(ctx, "invalid refresh token", "error", err)
		return nil, common.NewDomainError(common.ErrUnauthorized, "invalid_refresh_token", "invalid refresh token")
	}

	// Find refresh token in database
	storedToken, err := uc.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil || storedToken == nil {
		uc.logger.WarnContext(ctx, "refresh token not found in database", "user_id", userID)
		return nil, common.NewDomainError(common.ErrUnauthorized, "invalid_refresh_token", "refresh token not found")
	}

	// Check if token is expired
	if storedToken.IsExpired() {
		uc.logger.WarnContext(ctx, "expired refresh token", "user_id", userID)
		uc.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash) // Clean up expired token
		return nil, common.NewDomainError(common.ErrUnauthorized, "expired_refresh_token", "refresh token expired")
	}

	// Get user
	foundUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || foundUser == nil {
		uc.logger.ErrorContext(ctx, "user not found during token refresh", "user_id", userID)
		return nil, common.NewDomainError(common.ErrUnauthorized, "user_not_found", "user not found")
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := uc.tokenService.GenerateTokens(foundUser.ID, foundUser.Email.String())
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to generate new tokens", "user_id", userID, "error", err)
		return nil, common.NewDomainError(common.ErrInternal, "token_generation_failed", "failed to generate new tokens")
	}

	// Delete old refresh token and save new one
	if err := uc.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash); err != nil {
		uc.logger.WarnContext(ctx, "failed to delete old refresh token", "user_id", userID, "error", err)
	}

	if err := uc.saveRefreshToken(ctx, foundUser.ID, newRefreshToken); err != nil {
		uc.logger.ErrorContext(ctx, "failed to save new refresh token", "user_id", userID, "error", err)
		return nil, err
	}

	uc.logger.InfoContext(ctx, "token refresh successful", "user_id", userID)

	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
		User:         foundUser,
	}, nil
}

// GetProfile retrieves user profile by ID
func (uc *UseCase) GetProfile(ctx context.Context, userID common.ID) (*user.User, error) {
	uc.logger.InfoContext(ctx, "getting user profile", "user_id", userID)

	foundUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || foundUser == nil {
		uc.logger.ErrorContext(ctx, "user not found", "user_id", userID)
		return nil, common.NewDomainError(common.ErrNotFound, "user_not_found", "user not found")
	}

	return foundUser, nil
}

// saveRefreshToken saves a refresh token to the database
func (uc *UseCase) saveRefreshToken(ctx context.Context, userID common.ID, refreshToken string) error {
	// Extract token hash from refresh token (this depends on your token service implementation)
	_, tokenHash, err := uc.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return common.NewDomainError(common.ErrInternal, "invalid_token", "invalid refresh token format")
	}

	refreshTokenEntity := &user.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
	}

	if err := uc.refreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return common.NewDomainErrorWithCause(err, "token_save_failed", "failed to save refresh token")
	}

	return nil
}