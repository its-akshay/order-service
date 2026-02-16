package service

import (
	"context"
	"log"
	"time"
)

type Job struct {
	OrderID string
}

type WorkerPool struct {
	jobQueue chan Job
	workers  int
	service  *OrderService
}

func NewWorkerPool(workers int, service *OrderService) *WorkerPool {
	return &WorkerPool{
		jobQueue: make(chan Job, 100),
		workers:  workers,
		service:  service,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	for {
		select {
		case job := <-wp.jobQueue:
			log.Printf("Worker %d processing order %s", id, job.OrderID)
			wp.process(job)
		case <-ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) process(job Job) {
	// simulate payment delay
	time.Sleep(2 * time.Second)

	// update order status
	wp.service.UpdateOrderStatus(context.Background(), job.OrderID, "completed")
}
