Ainyx User API

A RESTful API built using Go (Fiber) to manage users with their name and DOB. It dynamically calculates a user’s age.

Project Type: Intern assignment for Ainyx

Features

CRUD operations for users

Age calculated from DOB

Input validation & proper error handling

Tech Stack

Backend: Go, Fiber

Database: PostgreSQL

SQL Queries: SQLC

Validation: go-playground/validator

API Endpoints
Method	Endpoint	Description
POST	/users	Create a user
GET	/users/:id	Get user by ID
PUT	/users/:id	Update user by ID
DELETE	/users/:id	Delete user by ID
GET	/users	List all users

Example Create Request:

POST /users
{
  "name": "Alice",
  "dob": "1990-05-10"
}


Response:

{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}

Setup & Run

Clone the repo:

git clone https://github.com/PagadalaAdarsh/ainyx-user-api.git
cd ainyx-user-api


Install dependencies:

go mod download


Set .env file with DB config:

DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_db_name
PORT=8080


Run server:

go run ./cmd/server/main.go


Server runs at: http://localhost:8080

Author: Pagadala Adarsh
GitHub: https://github.com/PagadalaAdarsh
LinkedIn: https://www.linkedin.com/in/adarsh-pagadala-63ba7227b
