package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Himany/gofermart/internal/logger"
	"go.uber.org/zap"
)

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	/*
		ТОЛЬКО ДЛЯ АВТОРИЗОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ

		Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых новых к самым старым.
		Формат даты — RFC3339.

		Доступные статусы обработки расчётов:
		NEW — заказ загружен в систему, но не попал в обработку;
		PROCESSING — вознаграждение за заказ рассчитывается;
		INVALID — система расчёта вознаграждений отказала в расчёте;
		PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.

		200 — успешная обработка запроса.
			Формат ответа:
			[
				{
					"number": "9278923470",
					"status": "PROCESSED",
					"accrual": 500,
					"uploaded_at": "2020-12-10T15:15:45+03:00"
				},
				{
					"number": "12345678903",
					"status": "PROCESSING",
					"uploaded_at": "2020-12-10T15:12:01+03:00"
				},
				{
					"number": "346436439",
					"status": "INVALID",
					"uploaded_at": "2020-12-09T16:09:53+03:00"
				}
			]
		204 — нет данных для ответа.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
	userID, isAuth := h.authFromRequest(r)
	if userID <= 0 || !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.Repo.ListUserOrders(userID)
	if err != nil {
		logger.Log.Error("GetOrders (ListUserOrders)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	payload, err := json.Marshal(orders)
	if err != nil {
		logger.Log.Error("GetOrders (Marshal)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logger.Log.Warn("GetOrders (write)", zap.Error(err))
		return
	}
}
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	/*
		ТОЛЬКО ДЛЯ АВТОРИЗОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ

		200 — успешная обработка запроса.
			Формат ответа:
			{
				"current": 500.5,
				"withdrawn": 42
			}
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
	userID, isAuth := h.authFromRequest(r)
	if userID <= 0 || !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := h.Repo.GetBalanceParts(userID)
	if err != nil {
		logger.Log.Error("GetBalance (GetBalanceParts)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		Current:   current,
		Withdrawn: withdrawn,
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		logger.Log.Error("GetBalance (Marshal)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logger.Log.Warn("GetBalance (write)", zap.Error(err))
		return
	}
}

func (h *Handler) GetBonusWithdrawals(w http.ResponseWriter, r *http.Request) {
	/*
		ТОЛЬКО ДЛЯ АВТОРИЗОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ

		Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым.
		Формат даты — RFC3339.

		200 — успешная обработка запроса.
			Формат ответа:
			[
				{
					"order": "2377225624",
					"sum": 500,
					"processed_at": "2020-12-09T16:09:57+03:00"
				}
			]
		204 — нет ни одного списания.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
	userID, isAuth := h.authFromRequest(r)
	if userID <= 0 || !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.Repo.ListWithdrawals(userID)
	if err != nil {
		logger.Log.Error("GetBonusWithdrawals (ListWithdrawals)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	payload, err := json.Marshal(withdrawals)
	if err != nil {
		logger.Log.Error("GetBonusWithdrawals (Marshal)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		logger.Log.Warn("GetBonusWithdrawals (write)", zap.Error(err))
		return
	}
}
