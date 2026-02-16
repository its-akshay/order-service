package service

import (
	"context"
	"encoding/json"
	"order-service/internal/cache"
	"order-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	Redis *cache.RedisClient
}

func NewOrderService(redis *cache.RedisClient) *OrderService {
	return &OrderService{Redis: redis}
}

func (s *OrderService) CreateOrder(ctx context.Context, amount float64) (*models.Order, error) {

	order := &models.Order{
		ID:        uuid.New().String(),
		Amount:    amount,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	err = s.Redis.Set(ctx, "order:"+order.ID, string(data))
	if err != nil {
		return nil, err
	}

	return order, nil
}
