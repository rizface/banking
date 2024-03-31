package handlers

import (
	"banking/api/responses"
	"banking/db"
	"banking/db/entity"
	"banking/internal/utils"
	"errors"
	"net/http"
	"net/mail"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Database *db.User
}

func validateUser(req struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}) error {
	lenEmail := len(req.Email)
	lenPassword := len(req.Password)
	lenName := len(req.Name)

	if lenEmail == 0 || lenPassword == 0 || lenName == 0 {
		return errors.New("email and password are required")
	}

	// validate email
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errors.New("invalid email format")
	}

	// validate name

	if lenName < 5 {
		return errors.New("name length must be at least 5 characters")
	}
	if lenName > 50 {
		return errors.New("name length cannot exceed 50 characters")
	}

	// validate password

	if lenPassword < 5 {
		return errors.New("password length must be at least 5 characters")
	}
	if lenPassword > 15 {
		return errors.New("password length cannot exceed 15 characters")
	}

	return nil
}

func validateLogin(req struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}) error {
	lenEmail := len(req.Email)
	lenPassword := len(req.Password)

	if lenEmail == 0 || lenPassword == 0 {
		return errors.New("email and password are required")
	}

	// validate email
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errors.New("invalid email format")
	}

	// validate password

	if lenPassword < 5 {
		return errors.New("password length must be at least 5 characters")
	}
	if lenPassword > 15 {
		return errors.New("password length cannot exceed 15 characters")
	}

	return nil
}

func (u *User) Register(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Validate request body
	if err := validateUser(req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	// Create user object
	usr := entity.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	// Register user
	result, err := u.Database.Register(ctx.UserContext(), usr)
	if err != nil {
		if err.Error() == "EXISTING_EMAIL" {
			return responses.ErrorConflict(ctx, err.Error())
		}

		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.Email, result.Id)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"name":        result.Name,
			"email":       result.Email,
			"accessToken": accessToken,
		},
	})
}

func (u *User) Login(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := validateLogin(req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	// login user
	result, err := u.Database.Login(ctx.UserContext(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "USER_NOT_FOUND" {
			return responses.ErrorNotFound(ctx, err.Error())
		}

		if err.Error() == "INVALID_PASSWORD" {
			return responses.ErrorBadRequest(ctx, err.Error())
		}

		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.Email, result.Id)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"name":        result.Name,
			"email":       result.Email,
			"accessToken": accessToken,
		},
	})
}
