package models

import "time"

type User struct {
	ID   int       `json:"id"`
	Name string    `json:"name" validate:"required"` // Name is required
	DOB  time.Time `json:"dob" validate:"required"`  // DOB is required
	Age  int       `json:"age,omitempty"`            // Calculated dynamically
}
