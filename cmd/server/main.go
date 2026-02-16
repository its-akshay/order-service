package main

import (
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/handler"
	"order-service/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := echo.New()

	redisClient := cache.NewRedisClient()
	orderService := service.NewOrderService(redisClient)
	orderHandler := handler.NewOrderHandler(orderService)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.POST("/orders", orderHandler.CreateOrder)
	
	e.GET("/orders/:id", orderHandler.GetOrder)


	e.Logger.Fatal(e.Start(":8080"))
}
