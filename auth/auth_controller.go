package auth

import (
	"github.com/hberkayozdemir/hypecoin-be/errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) Handler {
	return Handler{
		Service: service,
	}
}

func (h *Handler) AuthUserHandler(c *fiber.Ctx) error {
	log.Println("User authorization")
	bearerToken := c.Get("Authorization")
	if h.Service.VerifyToken(bearerToken, "user") || h.Service.VerifyToken(bearerToken, "admin") {
		return c.Next()
	}
	log.Println("Unauthorized user")
	return &fiber.Error{Code: 401, Message: errors.Unauthorized.Error()}
}

func (a *Handler) AuthAdminHandler(c *fiber.Ctx) error {
	log.Println("Admin authorization")
	bearerToken := c.Get("Authorization")
	if a.Service.VerifyToken(bearerToken, "admin") {
		return c.Next()
	}
	log.Println("Unauthorized admin")
	return &fiber.Error{Code: 401, Message: errors.Unauthorized.Error()}
}

func (a *Handler) AuthEditorHandler(c *fiber.Ctx) error {
	log.Println("Editor authorization")
	bearerToken := c.Get("Authorization")
	if a.Service.VerifyToken(bearerToken, "editor") {
		return c.Next()
	}
	log.Println("Unauthorized editor")
	return &fiber.Error{Code: 401, Message: errors.Unauthorized.Error()}
}
