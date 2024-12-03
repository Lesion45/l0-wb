package v1

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"wb-internship-l0/internal/service"
	"wb-internship-l0/pkg/validation"
)

type orderRoutes struct {
	log          *zap.Logger
	orderService service.Order
}

func newRoutes(log *zap.Logger, g *fiber.Router, orderService service.Order) {
	r := orderRoutes{
		log:          log,
		orderService: orderService,
	}

	(*g).Get("/get_order", r.getOrder)
}

type Request struct {
	ID string `json:"order_uid" validate:"required"`
}

func (r *orderRoutes) getOrder(c *fiber.Ctx) error {
	const op = "v1.orderRoutes.getOrder"

	var req Request

	if err := c.BodyParser(&req); err != nil {
		r.log.Error("failed to decode request body",
			zap.String("op", op),
			zap.String("route", "api/v1/get_order"),
			zap.Error(err),
		)

		return c.SendStatus(fiber.StatusBadRequest)
	}

	r.log.Info("request body decoded")
	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		r.log.Error("invalid request",
			zap.String("op", op),
			zap.Error(err),
		)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validation.ValidataionError(validateErr),
		})
	}

	data, err := r.orderService.GetOrder(context.Background(), req.ID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			r.log.Warn("order not found",
				zap.String("op", op),
				zap.String("route", "api/v1/get_order"),
				zap.String("orderID", req.ID),
				zap.Error(err),
			)

			return c.SendStatus(fiber.StatusNotFound)
		}

		r.log.Error("failed to get order",
			zap.String("op", op),
			zap.String("route", "api/v1/get_order"),
			zap.String("orderID", req.ID),
			zap.Error(err),
		)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(data)
}
