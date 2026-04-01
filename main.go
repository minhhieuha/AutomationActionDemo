package main

import (
	"log"

	"testing-demo/delivery/cronjob"
	httpDelivery "testing-demo/delivery/http"
	"testing-demo/repository"
	"testing-demo/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// === Repository Layer ===
	orderRepo := repository.NewInMemoryOrderRepo()

	// === Usecase Layer ===
	orderUC := usecase.NewOrderUsecase(orderRepo)

	// === Cronjob Service ===
	deliveryCron := cronjob.NewDeliveryCron(orderUC)
	deliveryCron.Start()
	defer deliveryCron.Stop()

	// === HTTP Delivery (Gin) ===
	router := gin.Default()
	httpDelivery.NewOrderHandler(router, orderUC)

	log.Println("[SERVER] Starting HTTP server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("[SERVER] Failed to start: %v", err)
	}
}
