package handlers

import (
	"net/http"
	"strconv"
	"test_nanimai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	svc service.Balance
}

func NewBalanceHandler(svc service.Balance) *BalanceHandler {
	return &BalanceHandler{svc: svc}
}

// UpdateLimit godoc
// @Summary Обновляет лимит счёта
// @Description Увеличивает/уменьшает максимальный лимит по счёту
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path int true "ID счёта"
// @Param input body service.UpdateLimitInput true "Изменение лимита"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts/{account_id}/limit [put]
func (h *BalanceHandler) UpdateLimit(c *gin.Context) {
	accountID, _ := strconv.ParseInt(c.Param("account_id"), 10, 64)
	var input service.UpdateLimitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.AccountID = accountID
	if err := h.svc.UpdateLimit(c.Request.Context(), input.AccountID, input.Delta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// UpdateBalance godoc
// @Summary Изменяет баланс счёта
// @Description Изменяет текущий баланс счёта на указанную величину
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path int true "ID счёта"
// @Param input body service.UpdateBalanceInput true "Изменение баланса"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts/{account_id}/balance [put]
func (h *BalanceHandler) UpdateBalance(c *gin.Context) {
	accountID, _ := strconv.ParseInt(c.Param("account_id"), 10, 64)
	var input service.UpdateBalanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.AccountID = accountID
	if err := h.svc.UpdateBalance(c.Request.Context(), input.AccountID, input.Delta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// OpenReservation godoc
// @Summary Открывает резерв средств
// @Description Создаёт резерв на сумму на указанном счёте
// @Tags reservations
// @Accept json
// @Produce json
// @Param account_id path int true "ID счёта"
// @Param input body service.OpenReservationInput true "Параметры резерва"
// @Success 200 {object} service.ReservationDTO
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts/{account_id}/reservation [post]
func (h *BalanceHandler) OpenReservation(c *gin.Context) {
	accountID, _ := strconv.ParseInt(c.Param("account_id"), 10, 64)
	var input service.OpenReservationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.AccountID = accountID
	res, err := h.svc.OpenReservation(c.Request.Context(), input.OwnerServiceID, input.AccountID, input.Amount, input.IdempotencyKey, input.Timeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ConfirmReservation godoc
// @Summary Подтверждает резерв
// @Description Подтверждает ранее открытый резерв
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation_id path int true "ID резерва"
// @Param X-Owner-Service-ID header int true "ID сервиса-владельца"
// @Success 200 {string} string "OK"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reservations/{reservation_id}/confirm [post]
func (h *BalanceHandler) ConfirmReservation(c *gin.Context) {
	reservationID, _ := strconv.ParseInt(c.Param("reservation_id"), 10, 64)
	ownerID, _ := strconv.ParseInt(c.GetHeader("X-Owner-Service-ID"), 10, 64)
	if err := h.svc.ConfirmReservation(c.Request.Context(), reservationID, ownerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// CancelReservation godoc
// @Summary Отменяет резерв
// @Description Отменяет ранее открытый резерв
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation_id path int true "ID резерва"
// @Param X-Owner-Service-ID header int true "ID сервиса-владельца"
// @Success 200 {string} string "OK"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /reservations/{reservation_id}/cancel [post]
func (h *BalanceHandler) CancelReservation(c *gin.Context) {
	reservationID, _ := strconv.ParseInt(c.Param("reservation_id"), 10, 64)
	ownerID, _ := strconv.ParseInt(c.GetHeader("X-Owner-Service-ID"), 10, 64)
	if err := h.svc.CancelReservation(c.Request.Context(), reservationID, ownerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
