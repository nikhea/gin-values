package services

import (
	"gin-learn/dto"
	"gin-learn/models"
	"time"

	"github.com/google/uuid"
)


var Users []models.User



func CreateUser(
	req dto.CreateUserRequest,
) models.User {


	user := models.User{

		ID: uuid.New().String(),

		Name:req.Name,

		Email:req.Email,

		Age:req.Age,

		CreatedAt:time.Now(),

		UpdatedAt:time.Now(),
	}


	Users = append(Users,user)


	return user
}



func GetUsers() []models.User {

	return Users
}

func GetUserByID(id string) (*models.User, bool) {

	for _, user := range Users {

		if user.ID == id {

			return &user, true
		}
	}

	return nil, false
}

func UpdateUser(id string, req dto.CreateUserRequest) (*models.User, bool) {

	for i, user := range Users {

		if user.ID == id {

			Users[i].Name = req.Name

			Users[i].Email = req.Email

			Users[i].Age = req.Age

			Users[i].UpdatedAt = time.Now()

			return &Users[i], true
		}
	}

	return nil, false
}

func DeleteUser(id string) bool {

	for i, user := range Users {

		if user.ID == id {

			Users = append(Users[:i], Users[i+1:]...)

			return true
		}
	}

	return false
}