package services

import (
	"context"
	"log/slog"

	"grip/dto"
	"grip/models"
	"grip/repository"
	"grip/utils"

	"github.com/google/uuid"
)

func CreateUser(req dto.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		ID:    uuid.New().String(),
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsers() ([]models.User, error) {
	return repository.GetUsers()
}

// ListUsers returns a filtered page of users plus the total count.
func ListUsers(filter dto.UserFilter) ([]models.User, int64, error) {
	filter.Normalize()
	return repository.ListUsers(filter)
}

func GetUserByID(id string) (*models.User, error) {
	// Cache-aside: serve from Redis when available, fall back to Postgres.
	ctx := context.Background()
	var cached models.User
	if hit, err := utils.GetJSON(ctx, utils.UserKey(id), &cached); err != nil {
		slog.Debug("user cache miss (error)", "id", id, "error", err)
	} else if hit {
		return &cached, nil
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if err := utils.SetJSON(ctx, utils.UserKey(id), user, utils.DefaultCacheTTL); err != nil {
		slog.Debug("user cache set failed", "id", id, "error", err)
	}
	return user, nil
}

func UpdateUser(id string, req dto.CreateUserRequest) (*models.User, error) {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	if err := repository.UpdateUser(user); err != nil {
		return nil, err
	}

	// Invalidate the cached copy (best-effort: TTL bounds staleness anyway).
	if err := utils.Del(context.Background(), utils.UserKey(id)); err != nil {
		slog.Debug("user cache invalidate failed", "id", id, "error", err)
	}

	return user, nil
}

func DeleteUser(id string) error {
	// Ensure user exists so callers can return 404
	if _, err := repository.GetUserByID(id); err != nil {
		return err
	}

	if err := repository.DeleteUser(id); err != nil {
		return err
	}

	if err := utils.Del(context.Background(), utils.UserKey(id)); err != nil {
		slog.Debug("user cache invalidate failed", "id", id, "error", err)
	}
	return nil
}
