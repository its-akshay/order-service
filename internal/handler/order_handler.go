package handler

import (
	"net/http"
	"order-service/internal/service"

	"github.com/labstack/echo/v4"
)

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	Amount float64 `json:"amount" example:"99.99"`
}

// OrderHandler handles order-related endpoints
type OrderHandler struct {
	Service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{Service: s}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the specified amount
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order amount"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {

	var req CreateOrderRequest

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

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get details of an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c echo.Context) error {

	id := c.Param("id")

	order, err := h.Service.GetOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	if order == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "order not found",
		})
	}

	return c.JSON(http.StatusOK, order)
}
