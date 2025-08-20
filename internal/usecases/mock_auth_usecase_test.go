package usecases

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"auth-service/internal/models"
// 	"auth-service/internal/repositorys"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/google/uuid"
// 	"github.com/sirupsen/logrus"
// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func setupTest() (*authUseCase, *repositorys.MockUserRepository) {
// 	log := logrus.New()
// 	validate := validator.New()
// 	v := viper.New()
// 	v.Set("jwt.secret", "test-secret")

// 	mockRepo := new(repositorys.MockUserRepository)

// 	// NewAuthUseCase return interface AuthUseCase
// 	uc := NewAuthUseCase(mockRepo, log, validate, v).(*authUseCase)

// 	return uc, mockRepo
// }

// func TestSignup(t *testing.T) {
// 	t.Run("Success", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		user := &models.User{ID: uuid.New(), Email: "test@example.com", Status: "active", EmailVerified: true}
// 		profile := &models.UserProfile{FullName: "Test User"}
// 		security := &models.UserSecurity{Password: "$2a$10$..."}
// 		role := &models.ApplicationRole{Role: "user"}

// 		mockRepo.On("CreateUser", user, profile, security, role).Return(nil)

// 		ctx := context.Background()
// 		result, err := usecase.Signup(ctx, "test@example.com", "password123", "Test User")
// 		assert.NoError(t, err)
// 		assert.NotNil(t, result)
// 		assert.Equal(t, user.ID, result.ID)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("ValidationError", func(t *testing.T) {
// 		usecase, _ := setupTest()
// 		ctx := context.Background()
// 		_, err := usecase.Signup(ctx, "", "short", "")
// 		assert.Error(t, err)
// 	})
// }

// func TestSignin(t *testing.T) {
// 	t.Run("Success", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		user := &models.User{ID: uuid.New(), Email: "test@example.com", Status: "active", EmailVerified: true}
// 		security := &models.UserSecurity{Password: "$2a$10$..."} // Hashed password

// 		mockRepo.On("FindUserByEmail", "test@example.com").Return(user, nil)
// 		mockRepo.On("FindUserSecurityByUserID", user.ID).Return(security, nil)
// 		mockRepo.On("CreateRefreshToken", mock.Anything).Return(nil)

// 		ctx := context.Background()
// 		accessToken, refreshToken, err := usecase.Signin(ctx, "test@example.com", "password123")
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, accessToken)
// 		assert.NotEmpty(t, refreshToken)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("InvalidCredentials", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		user := &models.User{ID: uuid.New(), Email: "test@example.com"}
// 		security := &models.UserSecurity{Password: "$2a$10$invalidhash"}

// 		mockRepo.On("FindUserByEmail", "test@example.com").Return(user, nil)
// 		mockRepo.On("FindUserSecurityByUserID", user.ID).Return(security, nil)

// 		ctx := context.Background()
// 		_, _, err := usecase.Signin(ctx, "test@example.com", "wrongpass")
// 		assert.Error(t, err)
// 	})
// }

// func TestChangePassword(t *testing.T) {
// 	t.Run("Success", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		userID := uuid.New()
// 		security := &models.UserSecurity{Password: "$2a$10$oldhash"}

// 		mockRepo.On("FindUserSecurityByUserID", userID).Return(security, nil)
// 		mockRepo.On("UpdateUserSecurity", userID, "$2a$10$newhash").Return(nil)

// 		ctx := context.Background()
// 		err := usecase.ChangePassword(ctx, userID, "oldpass", "newpass123")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("InvalidOldPassword", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		userID := uuid.New()
// 		security := &models.UserSecurity{Password: "$2a$10$oldhash"}

// 		mockRepo.On("FindUserSecurityByUserID", userID).Return(security, nil)

// 		ctx := context.Background()
// 		err := usecase.ChangePassword(ctx, userID, "wrongpass", "newpass123")
// 		assert.Error(t, err)
// 	})
// }

// func TestRefreshToken(t *testing.T) {
// 	t.Run("Success", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		token := &models.RefreshToken{
// 			SourceUserID: uuid.New(),
// 			TokenHash:    "validtoken",
// 			ExpiresAt:    time.Now().Add(24 * time.Hour),
// 		}
// 		user := &models.User{ID: token.SourceUserID, Email: "test@example.com"}

// 		mockRepo.On("FindRefreshToken", "validtoken").Return(token, nil)
// 		mockRepo.On("FindUserByEmail", token.SourceUserID.String()).Return(user, nil)
// 		mockRepo.On("CreateRefreshToken", mock.Anything).Return(nil)

// 		ctx := context.Background()
// 		newToken, err := usecase.RefreshToken(ctx, "validtoken")
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, newToken)
// 	})

// 	t.Run("InvalidToken", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		mockRepo.On("FindRefreshToken", "invalidtoken").Return(nil, assert.AnError)

// 		ctx := context.Background()
// 		_, err := usecase.RefreshToken(ctx, "invalidtoken")
// 		assert.Error(t, err)
// 	})
// }

// func TestChangeRole(t *testing.T) {
// 	t.Run("Success", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		userID := uuid.New()

// 		mockRepo.On("AssignRole", userID, "admin").Return(nil)

// 		ctx := context.Background()
// 		err := usecase.ChangeRole(ctx, userID, "admin")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("Failure", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()
// 		userID := uuid.New()

// 		mockRepo.On("AssignRole", userID, "invalid").Return(assert.AnError)

// 		ctx := context.Background()
// 		err := usecase.ChangeRole(ctx, userID, "invalid")
// 		assert.Error(t, err)
// 	})
// }

// func TestSignout(t *testing.T) {
// 	t.Run("Success", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()

// 		mockRepo.On("RevokeRefreshToken", "validtoken").Return(nil)

// 		ctx := context.Background()
// 		err := usecase.Signout(ctx, "validtoken")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("Failure", func(t *testing.T) {
// 		usecase, mockRepo := setupTest()

// 		mockRepo.On("RevokeRefreshToken", "invalidtoken").Return(assert.AnError)

// 		ctx := context.Background()
// 		err := usecase.Signout(ctx, "invalidtoken")
// 		assert.Error(t, err)
// 	})
// }
