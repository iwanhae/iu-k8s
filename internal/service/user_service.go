package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"iu-k8s.linecorp.com/server/internal/api"
	"iu-k8s.linecorp.com/server/internal/repository"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetHealth returns the health status of the application
func (s *UserService) GetHealth(ctx context.Context) (*api.HealthResponse, error) {
	now := time.Now()
	status := api.Healthy

	// Here you could add actual health checks (database, external services, etc.)

	return &api.HealthResponse{
		Status:    status,
		Timestamp: now,
		Version:   stringPtr("1.0.0"),
	}, nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) (*api.UsersResponse, error) {
	users, total, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	apiUsers := make([]api.User, len(users))
	for i, user := range users {
		apiUsers[i] = *userToAPI(&user)
	}

	hasMore := offset+len(users) < total

	return &api.UsersResponse{
		Users: apiUsers,
		Pagination: api.PaginationInfo{
			Limit:   limit,
			Offset:  offset,
			Total:   total,
			HasMore: &hasMore,
		},
	}, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req api.CreateUserRequest) (*api.User, error) {
	user := repository.User{
		ID:        uuid.New().String(),
		Email:     string(req.Email),
		Name:      req.Name,
		Avatar:    req.Avatar,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userToAPI(&user), nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*api.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return userToAPI(user), nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, userID string, req api.UpdateUserRequest) (*api.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Email != nil {
		user.Email = string(*req.Email)
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return userToAPI(user), nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Helper function to convert repository.User to api.User
func userToAPI(user *repository.User) *api.User {
	// Parse the UUID from string
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		// fallback to zero UUID if parsing fails
		userUUID = uuid.UUID{}
	}

	return &api.User{
		Id:        types.UUID(userUUID),
		Email:     types.Email(user.Email),
		Name:      user.Name,
		Avatar:    user.Avatar,
		IsActive:  &user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
