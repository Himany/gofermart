package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Himany/gofermart/internal/logger"
	"github.com/Himany/gofermart/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	/*
		{
			"login": "<login>",
			"password": "<password>"
		}

		Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization

		200 — пользователь успешно зарегистрирован и аутентифицирован;
		400 — неверный формат запроса;
		409 — логин уже занят;
		500 — внутренняя ошибка сервера.
	*/

	var registerRequest models.AuthDataRequest
	var buf bytes.Buffer

	//читаем тело запроса
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("Register (ReadFrom)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//десериализуем JSON в Visitor
	if err = json.Unmarshal(buf.Bytes(), &registerRequest); err != nil {
		logger.Log.Error("Register (Unmarshal)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Проверка объекта
	if err = validateRegisterJSON(registerRequest); err != nil {
		logger.Log.Error("Register (validateRegisterJSON)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	login := strings.ToLower(registerRequest.Login)
	login = strings.TrimSpace(login)
	uid, err := h.Repo.CreateUser(login, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			return
		}
		logger.Log.Error("Register (CreateUser)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.signJWT(uid, registerRequest.Login)
	if err != nil {
		logger.Log.Error("Register (signJWT)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.setJWTHeader(w, token)
	h.setJWTCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	/*
		{
			"login": "<login>",
			"password": "<password>"
		}

		Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization

		200 — пользователь успешно аутентифицирован;
		400 — неверный формат запроса;
		401 — неверная пара логин/пароль;
		500 — внутренняя ошибка сервера.
	*/
	var loginRequest models.AuthDataRequest
	var buf bytes.Buffer

	//читаем тело запроса
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("Login (ReadFrom)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//десериализуем JSON в Visitor
	if err = json.Unmarshal(buf.Bytes(), &loginRequest); err != nil {
		logger.Log.Error("Login (Unmarshal)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Проверка объекта
	if err = validateRegisterJSON(loginRequest); err != nil {
		logger.Log.Error("Login (validateRegisterJSON)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, hash, err := h.Repo.GetUserByLogin(loginRequest.Login)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(loginRequest.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := h.signJWT(id, loginRequest.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.setJWTHeader(w, token)
	h.setJWTCookie(w, token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	/*
		ТОЛЬКО ДЛЯ АВТОРИЗОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ

		Номер заказа передаётся в теле запроса в формате plain/text.
		12345678903

		Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна. https://ru.wikipedia.org/wiki/Алгоритм_Луна

		200 — номер заказа уже был загружен этим пользователем;
		202 — новый номер заказа принят в обработку;
		400 — неверный формат запроса;
		401 — пользователь не аутентифицирован;
		409 — номер заказа уже был загружен другим пользователем;
		422 — неверный формат номера заказа;
		500 — внутренняя ошибка сервера.
	*/
	userID, isAuth := h.authFromRequest(r)
	if userID <= 0 || !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, 512))
	if err != nil {
		logger.Log.Error("AddOrder (ReadAll)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ownerID, found, err := h.Repo.GetOrderOwner(orderNumber)
	if err != nil {
		logger.Log.Error("AddOrder (GetOrderOwner)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if found {
		if ownerID == userID {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
		return
	}

	if err := h.Repo.AddOrder(userID, orderNumber); err != nil {
		logger.Log.Error("AddOrder (AddOrder)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if h.AccrualPool != nil {
		h.AccrualPool.Submit(userID, orderNumber)
	}
}

func (h *Handler) BonusWithdraw(w http.ResponseWriter, r *http.Request) {
	/*
		ТОЛЬКО ДЛЯ АВТОРИЗОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ

		Номер заказа представляет собой гипотетический номер нового заказа пользователя,
		в счёт оплаты которого списываются баллы.

		Примечание: для успешного списания достаточно успешной регистрации запроса,
		никаких внешних систем начисления не предусмотрено и не требуется реализовывать.

		Формат запроса:
		{
			"order": "2377225624",
			"sum": 751
		}
		Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.

		200 — успешная обработка запроса;
		401 — пользователь не авторизован;
		402 — на счету недостаточно средств;
		422 — неверный номер заказа;
		500 — внутренняя ошибка сервера.
	*/

	userID, isAuth := h.authFromRequest(r)
	if userID <= 0 || !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var withdrawRequest models.BonusWithdrawRequest
	var buf bytes.Buffer

	//читаем тело запроса
	defer r.Body.Close()
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Log.Error("BonusWithdraw (ReadFrom)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//десериализуем JSON в Visitor
	if err = json.Unmarshal(buf.Bytes(), &withdrawRequest); err != nil {
		logger.Log.Error("BonusWithdraw (Unmarshal)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Проверка объекта
	if err = validateWithdrawJSON(withdrawRequest); err != nil {
		logger.Log.Error("BonusWithdraw (validateWithdrawJSON)", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(withdrawRequest.Order)
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := h.Repo.AddWithdraw(userID, orderNumber, withdrawRequest.Sum); err != nil {
		if strings.Contains(err.Error(), "insufficient funds") {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		logger.Log.Error("BonusWithdraw (AddWithdraw)", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
