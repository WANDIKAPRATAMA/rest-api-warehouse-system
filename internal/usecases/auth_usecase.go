package usecases

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"auth-service/internal/models"
	"auth-service/internal/repositorys"
	"auth-service/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Signup(ctx context.Context, email, password, fullName string) (*models.User, error)
	Signin(ctx context.Context, email, password string, deviceID *string) (string, string, error) // accessToken, refreshToken
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string, deviceID string) (string, string, error) // newAccessToken
	ChangeRole(ctx context.Context, userID uuid.UUID, role string) error
	Signout(ctx context.Context, tokenHash string) error
}

type authUseCase struct {
	repo     repositorys.UserRepository
	validate *validator.Validate
	log      *logrus.Logger
	config   *viper.Viper
}

func NewAuthUseCase(repo repositorys.UserRepository, log *logrus.Logger, validate *validator.Validate, config *viper.Viper) AuthUseCase {
	return &authUseCase{repo: repo, log: log, validate: validate, config: config}
}

func (u *authUseCase) Signup(ctx context.Context, email, password, fullName string) (*models.User, error) {

	if err := u.validate.Struct(&struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
		FullName string `validate:"required"`
	}{Email: email, Password: password, FullName: fullName}); err != nil {
		return nil, err
	}

	exist, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{Email: email, Status: "active", EmailVerified: true}
	profile := &models.UserProfile{FullName: fullName}
	security := &models.UserSecurity{Password: string(hashedPassword)}
	role := &models.ApplicationRole{Role: "user"}

	if err := u.repo.CreateUser(user, profile, security, role); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *authUseCase) Signin(ctx context.Context, email, password string, deviceID *string) (string, string, error) {
	user, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	security, err := u.repo.FindUserSecurityByUserID(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(security.Password), []byte(password)); err != nil {
		return "", "", err
	}

	role, err := u.repo.FindUserRoleByUserID(user.ID)
	if err != nil {
		return "", "", err
	}

	secret := u.config.GetString("jwt.secret")

	accessToken, err := utils.GenerateAccessToken(secret, user.ID, user.Email, role)
	if err != nil {
		return "", "", err
	}

	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	refresh := &models.RefreshToken{
		SourceUserID: user.ID,
		TokenHash:    refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		DeviceID:     *deviceID,
	}
	if err := u.repo.CreateRefreshToken(refresh); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *authUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	security, err := u.repo.FindUserSecurityByUserID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(security.Password), []byte(oldPassword)); err != nil {
		return err
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.UpdateUserSecurity(userID, string(hashedNewPassword))
}

func (u *authUseCase) RefreshToken(ctx context.Context, refreshToken string, deviceID string) (string, string, error) {

	u.log.Println("device id", deviceID, "refresh token", refreshToken)
	storedToken, err := u.repo.FindRefreshToken(refreshToken, deviceID)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	u.log.Println("stored token", storedToken)
	if storedToken.RevokedAt != nil && !storedToken.RevokedAt.IsZero() {
		return "", "", fmt.Errorf("refresh token revoked")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return "", "", fmt.Errorf("refresh token expired")
	}

	// Ambil user
	user, err := u.repo.FindUserByID(storedToken.SourceUserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	role, err := u.repo.FindUserRoleByUserID(user.ID)
	if err != nil {
		return "", "", err
	}

	secret := u.config.GetString("jwt.secret")

	// Generate access token baru
	accessToken, err := utils.GenerateAccessToken(secret, user.ID, user.Email, role)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token baru (opsional, best practice rotate)
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	newRefreshToken := hex.EncodeToString(refreshTokenBytes)

	storedToken.TokenHash = newRefreshToken
	storedToken.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	storedToken.LastUsedAt = time.Now()

	if err := u.repo.UpdateRefreshToken(storedToken); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (u *authUseCase) ChangeRole(ctx context.Context, userID uuid.UUID, role string) error {
	return u.repo.AssignRole(userID, role)
}

func (u *authUseCase) Signout(ctx context.Context, tokenHash string) error {
	return u.repo.RevokeRefreshToken(tokenHash)
}
