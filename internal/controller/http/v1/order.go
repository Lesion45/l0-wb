package v1

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"wb-internship-l0/internal/service"
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
	ID string `json:"order_uid"`
}

func (r *orderRoutes) getOrder(c *fiber.Ctx) error {
	var req Request

	if err := c.BodyParser(&req); err != nil {
		r.log.Error("failed to get order_id from request",
			zap.String("op", "v1.orderRoutes.getOrder"),
			zap.String("route", "api/v1/get_order"),
			zap.Error(err),
		)

		return c.SendStatus(fiber.StatusBadRequest)
	}

	data, err := r.orderService.GetOrder(context.Background(), req.ID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			r.log.Warn("order not found",
				zap.String("op", "v1.orderRoutes.getOrder"),
				zap.String("route", "api/v1/get_order"),
				zap.String("orderID", req.ID),
				zap.Error(err),
			)

			return c.SendStatus(fiber.StatusNotFound)
		}

		r.log.Error("failed to get order",
			zap.String("op", "v1.orderRoutes.getOrder"),
			zap.String("route", "api/v1/get_order"),
			zap.String("orderID", req.ID),
			zap.Error(err),
		)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(data)
}
