package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/oapi-codegen/runtime/types"
	"iu-k8s.linecorp.com/server/internal/api"
	"iu-k8s.linecorp.com/server/internal/service"
)

// UserHandler implements the ServerInterface from generated code
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetHealth handles the health check endpoint
func (h *UserHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health, err := h.userService.GetHealth(ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "HEALTH_CHECK_FAILED",
			Message: "Failed to get health status",
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, health)
}

// ListUsers handles listing users with pagination
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request, params api.ListUsersParams) {
	ctx := r.Context()

	// Set default values and parse query parameters
	limit := 10
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	users, err := h.userService.ListUsers(ctx, limit, offset)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "LIST_USERS_FAILED",
			Message: "Failed to list users",
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}

// CreateUser handles creating a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req api.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	user, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "CREATE_USER_FAILED",
			Message: "Failed to create user",
		})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// GetUserById handles getting a user by ID
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request, userId types.UUID) {
	ctx := r.Context()

	user, err := h.userService.GetUserByID(ctx, userId.String())
	if err != nil {
		if err.Error() == "user not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, api.ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "GET_USER_FAILED",
			Message: "Failed to get user",
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, userId types.UUID) {
	ctx := r.Context()

	var req api.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	user, err := h.userService.UpdateUser(ctx, userId.String(), req)
	if err != nil {
		if err.Error() == "failed to get user: user not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, api.ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "UPDATE_USER_FAILED",
			Message: "Failed to update user",
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, userId types.UUID) {
	ctx := r.Context()

	err := h.userService.DeleteUser(ctx, userId.String())
	if err != nil {
		if err.Error() == "failed to delete user: user not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, api.ErrorResponse{
				Error:   "USER_NOT_FOUND",
				Message: "User not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.ErrorResponse{
			Error:   "DELETE_USER_FAILED",
			Message: "Failed to delete user",
		})
		return
	}

	render.Status(r, http.StatusNoContent)
}
