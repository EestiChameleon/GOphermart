package handlers

import (
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/EestiChameleon/GOphermart/internal/ctxfunc"
	"io"
	"net/http"
)

// UserAddOrder загрузка пользователем номера заказа для расчёта;
/*
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
12345678903

Возможные коды ответа:
200 — номер заказа уже был загружен этим пользователем; +
202 — новый номер заказа принят в обработку; +
400 — неверный формат запроса; +
401 — пользователь не аутентифицирован; +
409 — номер заказа уже был загружен другим пользователем; +
422 — неверный формат номера заказа; +
500 — внутренняя ошибка сервера. +
*/
func UserAddOrder(w http.ResponseWriter, r *http.Request) {
	userID := ctxfunc.GetUserIDFromCTX(r.Context())
	if userID < 1 {
		resp.NoContent(w, http.StatusUnauthorized)
		return
	}

	// read body
	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	// check order number - empty or Luhn
	orderNumber := string(byteBody)
	if orderNumber == "" || !service.LuhnCheck(orderNumber) {
		resp.WriteString(w, http.StatusUnprocessableEntity, "invalid order number")
		return
	}

	cmlogger.Sug.Infow("User Add Order start", "UserID", userID, "Order Number", orderNumber)
	o := methods.NewOrder(userID, orderNumber)

	if err = o.CheckNumber(); err != nil {
		switch err {
		case methods.ErrOrderAlreadyExists:
			cmlogger.Sug.Infow("order already in process", "Number", orderNumber)
			resp.NoContent(w, http.StatusOK)
			return
		case methods.ErrOrderWrongOwner:
			cmlogger.Sug.Infow("new order conflict - owned by another user", "UserID", userID, "Number", orderNumber)
			resp.NoContent(w, http.StatusConflict)
		}
	}

	// after check is done - we are sure that the order is new = we can save it
	if err = o.Add(); err != nil { // new order add to table
		cmlogger.Sug.Error("order.Add err:%v", err)
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	cmlogger.Sug.Infow("new order accepted", "Number", orderNumber)
	resp.NoContent(w, http.StatusAccepted)
}
