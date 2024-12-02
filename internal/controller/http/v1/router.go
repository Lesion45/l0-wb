package v1

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"wb-internship-l0/internal/service"
)

func InitRouter(log *zap.Logger, app *fiber.App, services *service.Services) {
	v1 := app.Group("api/v1")

	newRoutes(log, &v1, services.Order)
}
