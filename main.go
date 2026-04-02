package main

import (
	"log"
	"os"

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

	// Render và các Cloud provider thường gán Port qua biến môi trường PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[SERVER] Starting HTTP server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[SERVER] Failed to start: %v", err)
	}
}
