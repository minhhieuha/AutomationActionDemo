package cronjob

import (
	"log"
	"testing-demo/domain"

	"github.com/robfig/cron/v3"
)

type DeliveryCron struct {
	orderUsecase domain.OrderUsecase
	scheduler    *cron.Cron
}

// NewDeliveryCron creates and starts delivery cronjob
func NewDeliveryCron(uc domain.OrderUsecase) *DeliveryCron {
	dc := &DeliveryCron{
		orderUsecase: uc,
		scheduler:    cron.New(cron.WithSeconds()),
	}

	// Run every 10 seconds: deliver pending orders
	_, err := dc.scheduler.AddFunc("*/10 * * * * *", dc.deliverOrders)
	if err != nil {
		log.Fatalf("[CRON] Failed to register delivery job: %v", err)
	}

	return dc
}

// Start begins the cron scheduler
func (dc *DeliveryCron) Start() {
	log.Println("[CRON] Delivery cronjob started (every 10 seconds)")
	dc.scheduler.Start()
}

// Stop gracefully stops the cron scheduler
func (dc *DeliveryCron) Stop() {
	log.Println("[CRON] Delivery cronjob stopped")
	dc.scheduler.Stop()
}

func (dc *DeliveryCron) deliverOrders() {
	log.Println("[CRON] Checking for pending orders...")
	err := dc.orderUsecase.DeliverPendingOrders()
	if err != nil {
		log.Printf("[CRON] Error delivering orders: %v\n", err)
		return
	}
	log.Println("[CRON] Pending orders delivered successfully")
}
