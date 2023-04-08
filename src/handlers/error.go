package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ERROR_INVALID_PAYLOAD  = &fiber.Error{Code: 400, Message: "Error. Invalid payload"}
	ERROR_CREATING_USER    = &fiber.Error{Code: 400, Message: "Error creating new user"}
	ERROR_GENERATING_JWT   = &fiber.Error{Code: 400, Message: "Please try again later"}
	ERROR_INVALID_EMAIL    = &fiber.Error{Code: 401, Message: "Error. Invalid email"}
	ERROR_INVALID_PASSWORD = &fiber.Error{Code: 401, Message: "Error. Invalid password"}
	ERROR_NOT_AUTHORIZED   = &fiber.Error{Code: 401, Message: "Error. Not authorized"}
	ERROR_UPDATING_DATA    = &fiber.Error{Code: 400, Message: "Error updating data"}
)

func ResponseWhenError(ctx *fiber.Ctx, err error) error {

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		err = ctx.Status(e.Code).JSON(e)
	}

	// Return a generic error in case a fiber.Error is not set
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Error{
				Code:    500,
				Message: "Internal Server Error",
			})
	}
	return nil
}
