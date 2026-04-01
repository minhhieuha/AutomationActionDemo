package domain

import "time"

// Order represents an order in the system
type Order struct {
	ID           int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	Item         string    `json:"item"`
	Status       string    `json:"status"` // "Pending", "Delivered"
	CreatedAt    time.Time `json:"created_at"`
}

// OrderRepository is the interface that the repository layer must implement
type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	GetPendingOrders() ([]*Order, error)
	UpdateStatus(id int, status string) error
}

// OrderUsecase is the interface that the usecase layer must implement
type OrderUsecase interface {
	CreateOrder(customerName, item string) (*Order, error)
	GetOrder(id int) (*Order, error)
	DeliverPendingOrders() error
}
