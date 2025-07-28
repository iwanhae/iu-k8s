package repository

import (
	"context"
	"errors"
	"sync"
	"time"
)

// User represents a user entity in the repository layer
type User struct {
	ID        string
	Email     string
	Name      string
	Avatar    *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	List(ctx context.Context, limit, offset int) ([]User, int, error)
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// InMemoryUserRepository is an in-memory implementation of UserRepository
// This is useful for development and testing
type InMemoryUserRepository struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

// List returns a paginated list of users
func (r *InMemoryUserRepository) List(ctx context.Context, limit, offset int) ([]User, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	total := len(r.users)
	users := make([]User, 0, limit)

	// Convert map to slice for pagination
	allUsers := make([]*User, 0, total)
	for _, user := range r.users {
		allUsers = append(allUsers, user)
	}

	// Apply pagination
	start := offset
	if start > len(allUsers) {
		start = len(allUsers)
	}

	end := start + limit
	if end > len(allUsers) {
		end = len(allUsers)
	}

	for i := start; i < end; i++ {
		users = append(users, *allUsers[i])
	}

	return users, total, nil
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	// Create a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return a copy to avoid external modifications
	userCopy := *user
	return &userCopy, nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Create a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}
