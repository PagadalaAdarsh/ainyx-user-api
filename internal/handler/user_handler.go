package handler

import (
	"strconv"

	db "github.com/adarsh/ainyx-task/db/sqlc"
	"github.com/adarsh/ainyx-task/internal/logger"
	"github.com/adarsh/ainyx-task/internal/models"
	"github.com/adarsh/ainyx-task/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}
	user, err := h.svc.Store().GetUserByID(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	logger.Log.Info("Fetched user", zap.Int("id", id))
	return c.JSON(h.svc.ToModel(user))
}

// POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	// Validate input
	if err := h.svc.ValidateUser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.svc.Store().CreateUser(c.Context(), db.CreateUserParams{
		Name: input.Name,
		Dob:  input.DOB,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
	}

	logger.Log.Info("Created user", zap.String("name", input.Name), zap.Time("dob", input.DOB))
	return c.Status(fiber.StatusCreated).JSON(h.svc.ToModel(user))
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	// Validate input
	if err := h.svc.ValidateUser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.svc.Store().UpdateUser(c.Context(), db.UpdateUserParams{
		ID:   int32(id),
		Name: input.Name,
		Dob:  input.DOB,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update user"})
	}

	logger.Log.Info("Updated user", zap.Int("id", id), zap.String("name", input.Name))
	return c.JSON(h.svc.ToModel(user))
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}

	if err := h.svc.Store().DeleteUser(c.Context(), int32(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete user"})
	}

	logger.Log.Info("Deleted user", zap.Int("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

// GET /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.svc.Store().ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch users"})
	}

	var result []models.User
	for _, u := range users {
		result = append(result, h.svc.ToModel(u))
	}

	logger.Log.Info("Listed users", zap.Int("count", len(result)))
	return c.JSON(result)
}
