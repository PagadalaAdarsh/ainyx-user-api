package service

import (
	"time"

	db "github.com/adarsh/ainyx-task/db/sqlc"
	"github.com/adarsh/ainyx-task/internal/models"
	"github.com/go-playground/validator/v10" // <- add this
)

// create a validator instance
var validate = validator.New()

type UserService struct {
	store *db.Queries
}

func NewUserService(store *db.Queries) *UserService {
	return &UserService{store: store}
}

// Getter to access store safely
func (s *UserService) Store() *db.Queries {
	return s.store
}

// CalculateAge calculates age from DOB
func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

// ToModel converts DB User to API model User with calculated age
func (s *UserService) ToModel(u db.User) models.User {
	return models.User{
		ID:   int(u.ID),
		Name: u.Name,
		DOB:  u.Dob,
		Age:  CalculateAge(u.Dob),
	}
}

// ✅ New function to validate User input
func (s *UserService) ValidateUser(u models.User) error {
	return validate.Struct(u)
}
