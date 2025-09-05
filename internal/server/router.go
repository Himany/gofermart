package server

import (
	"net/http"

	"github.com/Himany/gofermart/internal/logger"
	"github.com/Himany/gofermart/internal/middleware"
	"github.com/Himany/gofermart/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

func Router(handler *handlers.Handler, runAddress string) error {
	r := chi.NewRouter()

	r.Post("/api/user/register", middleware.CheckApplicationJSONContentType(handler.Register))              //регистрация пользователя;
	r.Post("/api/user/login", middleware.CheckApplicationJSONContentType(handler.Login))                    //аутентификация пользователя;
	r.Post("/api/user/orders", middleware.CheckPlainTextContentType(handler.AddOrder))                      //загрузка пользователем номера заказа для расчёта;
	r.Get("/api/user/orders", handler.GetOrders)                                                            //получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	r.Get("/api/user/balance", handler.GetBalance)                                                          //получение текущего баланса счёта баллов лояльности пользователя;
	r.Post("/api/user/balance/withdraw", middleware.CheckApplicationJSONContentType(handler.BonusWithdraw)) //запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	r.Get("/api/user/withdrawals", handler.GetBonusWithdrawals)                                             //получение информации о выводе средств с накопительного счёта пользователем.

	return http.ListenAndServe(runAddress, logger.RequestLogger(r))
}
