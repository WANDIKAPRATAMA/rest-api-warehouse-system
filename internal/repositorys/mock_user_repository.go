package repositorys

import (
	"auth-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User, profile *models.UserProfile, security *models.UserSecurity, role *models.ApplicationRole) error {
	args := m.Called(user, profile, security, role)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindUserSecurityByUserID(userID uuid.UUID) (*models.UserSecurity, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.UserSecurity), args.Error(1)
}

func (m *MockUserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockUserRepository) FindRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	args := m.Called(tokenHash)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockUserRepository) RevokeRefreshToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserSecurity(userID uuid.UUID, newPassword string) error {
	args := m.Called(userID, newPassword)
	return args.Error(0)
}

func (m *MockUserRepository) AssignRole(userID uuid.UUID, role string) error {
	args := m.Called(userID, role)
	return args.Error(0)
}
func (m *MockUserRepository) FindUserRoleByUserID(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
