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
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	data, _ := json.Marshal(order)

	if err := s.Redis.Set(ctx, "order:"+order.ID, string(data)); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {

	data, err := s.Redis.Get(ctx, "order:"+id)
	if err != nil {
		return nil, err
	}

	if data == "" {
		return nil, nil // not found
	}

	var order models.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id, status string) error {

	order, err := s.GetOrder(ctx, id)
	if err != nil || order == nil {
		return err
	}

	order.Status = status

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return s.Redis.Set(ctx, "order:"+id, string(data))
}
