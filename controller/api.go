package controller

import (
	"github.com/hberkayozdemir/hypecoin-be/errors"
	"github.com/hberkayozdemir/hypecoin-be/model"
	"github.com/hberkayozdemir/hypecoin-be/service"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	service *service.Service
}

func (api *API) SetupApp(app *fiber.App) {
	app.Post("/api/register", api.RegisterUserHandler)
	app.Post("/api/login", api.LoginUserHandler)
	app.Get("/api/activation/:userID", api.ActivationHandler)
	app.Post("/api/forgotPassword", api.ForgotPasswordHandler)
	app.Patch("/api/resetPassword/:userID", api.ResetPasswordHandler)
	app.Get("/api/users/:userID", api.GetUserHandler)
	app.Get("/news",api.GetNews)
}

func NewAPI(service *service.Service) API {
	return API{
		service: service,
	}
}

func (api *API) RegisterUserHandler(c *fiber.Ctx) error {
	userDTO := model.UserDTO{}

	err := c.BodyParser(&userDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := api.service.RegisterUser(userDTO)

	switch err {
	case nil:
		c.JSON(user)
		c.Status(fiber.StatusCreated)
	case errors.UserAlreadyRegistered:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) LoginUserHandler(c *fiber.Ctx) error {
	userCredentialsDTO := model.UserCredentialsDTO{}

	err := c.BodyParser(&userCredentialsDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	token, cookie, err := api.service.LoginUser(userCredentialsDTO)

	switch err {
	case nil:
		c.JSON(token)
		c.Cookie(cookie)
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusBadRequest)
	case errors.Unauthorized:
		c.Status(fiber.StatusUnauthorized)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ActivationHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	user, err := api.service.Activation(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ForgotPasswordHandler(c *fiber.Ctx) error {
	forgotPasswordDTO := model.ForgotPasswordDTO{}
	err := c.BodyParser(&forgotPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = api.service.ForgotPassword(forgotPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	case errors.UserNotActivated:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ResetPasswordHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	resetPasswordDTO := model.ResetPasswordDTO{}
	err := c.BodyParser(&resetPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = api.service.ResetPassword(userID, resetPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) GetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	if len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := api.service.GetUser(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}


func (api *API) GetNews(c *fiber.Ctx) error {
	

	
	

	news, err := api.service.GetNews()

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(news)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	updateUserDTO := model.UpdateUserDTO{}
	err := c.BodyParser(&updateUserDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	updatedUser, err := api.service.UpdateUser(userID, updateUserDTO)
	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(updatedUser)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

