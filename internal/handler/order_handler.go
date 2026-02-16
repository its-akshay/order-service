package handler

import (
	"net/http"
	"order-service/internal/service"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	Service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{Service: s}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	order, err := h.Service.CreateOrder(c.Request().Context(), req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create order",
		})
	}

	return c.JSON(http.StatusCreated, order)
}
