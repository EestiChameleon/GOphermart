package handlers

import (
	"errors"
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	db "github.com/EestiChameleon/GOphermart/internal/app/storage"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"io"
	"net/http"
)

// UserAddOrder загрузка пользователем номера заказа для расчёта;
/*
POST /api/user/orders HTTP/1.1
Content-Type: text/plain
12345678903

Возможные коды ответа:
200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.
*/
func UserAddOrder(w http.ResponseWriter, r *http.Request) {
	// read body
	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	// check if it's not empty
	orderNumber := string(byteBody)
	if orderNumber == "" {
		resp.WriteString(w, http.StatusBadRequest, "empty input")
		return
	}

	if !service.LuhnCheck(orderNumber) {
		resp.WriteString(w, http.StatusUnprocessableEntity, "invalid order number")
		return
	}

	cmlogger.Sug.Infow("User Add Order start", "UserID", db.Pool.ID, "Order Number", orderNumber)
	o := methods.NewOrder(orderNumber)
	if err = o.GetByNumber(); err != nil {
		if !errors.Is(err, db.ErrNotFound) { // in case NotFound - new order -> we can proceed
			cmlogger.Sug.Error("order.GetByNumber err:%v", err)
			resp.NoContent(w, http.StatusInternalServerError)
			return
		}

		if err = o.Add(); err != nil { // new order add to table
			cmlogger.Sug.Error("order.Add err:%v", err)
			resp.NoContent(w, http.StatusInternalServerError)
			return
		}

		cmlogger.Sug.Infow("new order accepted", "Number", orderNumber)
		resp.NoContent(w, http.StatusAccepted)
		return
	}

	if db.Pool.ID == o.UserID { // case found - we compare currentSession.userID and DB.orders.userID
		cmlogger.Sug.Infow("order already in process", "Number", orderNumber)
		resp.NoContent(w, http.StatusOK)
		return
	} else { //case userIDs don't match
		cmlogger.Sug.Infow("new order conflict", "UserID", db.Pool.ID, "Number", orderNumber)
		resp.NoContent(w, http.StatusConflict)
		return
	}
}
