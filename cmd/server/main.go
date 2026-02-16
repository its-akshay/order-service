package main

import (
	"context"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/handler"
	"order-service/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "order-service/docs" // Import generated docs
)

// @title Order Service API
// @version 1.0
// @description This is an order service API
// @host localhost:8080
// @BasePath /
func main() {
	e := echo.New()

	// Add CORS middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	redisClient := cache.NewRedisClient()
	orderService := service.NewOrderService(redisClient)
	orderHandler := handler.NewOrderHandler(orderService)
	ctx := context.Background()

	workerPool := service.NewWorkerPool(3, orderService)
	workerPool.Start(ctx)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Fix the swagger route
e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.POST("/orders", orderHandler.CreateOrder)

	e.GET("/orders/:id", orderHandler.GetOrder)

	e.Logger.Fatal(e.Start(":8080"))
}
