package rest

import (
	"github.com/gin-gonic/gin"
	handlers2 "test_nanimai/backend/internal/api/rest/handlers"
	"test_nanimai/backend/internal/service"
)

func RegisterRoutes(r *gin.Engine, svc service.Balance) {
	handler := handlers2.NewBalanceHandler(svc)

	r.PUT("/accounts/:account_id/limit", handler.UpdateLimit)
	r.PUT("/accounts/:account_id/balance", handler.UpdateBalance)
	r.POST("/accounts/:account_id/reservation", handler.OpenReservation)
	r.POST("/reservations/:reservation_id/confirm", handler.ConfirmReservation)
	r.POST("/reservations/:reservation_id/cancel", handler.CancelReservation)
}
