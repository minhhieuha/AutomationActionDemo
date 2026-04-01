package cronjob

import (
	"reflect"
	"testing"
	"testing-demo/domain"

	"github.com/robfig/cron/v3"
)

func TestNewDeliveryCron(t *testing.T) {
	type args struct {
		uc domain.OrderUsecase
	}
	tests := []struct {
		name string
		args args
		want *DeliveryCron
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeliveryCron(tt.args.uc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeliveryCron() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeliveryCron_Start(t *testing.T) {
	type fields struct {
		orderUsecase domain.OrderUsecase
		scheduler    *cron.Cron
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DeliveryCron{
				orderUsecase: tt.fields.orderUsecase,
				scheduler:    tt.fields.scheduler,
			}
			dc.Start()
		})
	}
}

func TestDeliveryCron_Stop(t *testing.T) {
	type fields struct {
		orderUsecase domain.OrderUsecase
		scheduler    *cron.Cron
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DeliveryCron{
				orderUsecase: tt.fields.orderUsecase,
				scheduler:    tt.fields.scheduler,
			}
			dc.Stop()
		})
	}
}

func TestDeliveryCron_deliverOrders(t *testing.T) {
	type fields struct {
		orderUsecase domain.OrderUsecase
		scheduler    *cron.Cron
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DeliveryCron{
				orderUsecase: tt.fields.orderUsecase,
				scheduler:    tt.fields.scheduler,
			}
			dc.deliverOrders()
		})
	}
}
