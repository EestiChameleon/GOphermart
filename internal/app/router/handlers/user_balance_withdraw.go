package handlers

import (
	"encoding/json"
	"errors"
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/EestiChameleon/GOphermart/internal/ctxfunc"
	"github.com/EestiChameleon/GOphermart/internal/models"
	"io"
	"net/http"
)

// UserBalanceWithdraw запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
/*
Формат запроса:
POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json
{
    "order": "2377225624",
    "sum": 751
}
Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.
Возможные коды ответа:
200 — успешная обработка запроса;
401 — пользователь не авторизован;
402 — на счету недостаточно средств;
422 — неверный номер заказа;
500 — внутренняя ошибка сервера.
*/

func UserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	userID := ctxfunc.GetUserIDFromCTX(r.Context())
	if userID < 1 {
		resp.NoContent(w, http.StatusUnauthorized)
		return
	}

	var b models.WithdrawData
	data, err := io.ReadAll(r.Body)
	if err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw read body err:%v", err)
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(data, &b); err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw Unmarshal body err:%v", err)
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	if b.Order == "" || !service.LuhnCheck(b.Order) {
		cmlogger.Sug.Errorf("UserBalanceWithdraw Empty Order or LuhnCheck err:%v", err)
		resp.NoContent(w, http.StatusUnprocessableEntity)
		return
	}

	// proceed new balance withdraw
	err = service.BalanceWithdraw(userID, b)
	if err != nil && !errors.Is(err, service.ErrWithdrawUnavailable) {
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, service.ErrWithdrawUnavailable) {
		resp.NoContent(w, http.StatusPaymentRequired)
		return
	}

	resp.NoContent(w, http.StatusOK)
}
