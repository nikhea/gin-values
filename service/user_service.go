package services

import (
	"gin-learn/dto"
	"gin-learn/models"
	"gin-learn/repository"

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

func GetUserByID(id string) (*models.User, error) {
	return repository.GetUserByID(id)
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

	return user, nil
}

func DeleteUser(id string) error {
	// Ensure user exists so callers can return 404
	if _, err := repository.GetUserByID(id); err != nil {
		return err
	}

	return repository.DeleteUser(id)
}
