package handlers

import (
	"encoding/json"
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
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
	var b models.WithdrawData
	data, err := io.ReadAll(r.Body)
	if err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw read body err:%v", err)
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &b)
	if err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw Unmarshal body err:%v", err)
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	if b.Order == "" || !service.LuhnCheck(b.Order) {
		cmlogger.Sug.Errorf("UserBalanceWithdraw Empty Order or LuhnCheck err:%v", err)
		resp.NoContent(w, http.StatusUnprocessableEntity)
		return
	}

	// get current balance and whole withdrawn
	blnc := methods.NewBalanceRecord()
	res, err := blnc.GetBalanceAndWithdrawnByUserID()
	if err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw GetBalanceAndWithdrawnByUserID err:%v", err)
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	if !res.Current.GreaterThanOrEqual(b.Sum.Decimal) {
		cmlogger.Sug.Errorf("UserBalanceWithdraw GreaterThanOrEqual %v ? %v", res.Current, b.Sum.Decimal)
		resp.NoContent(w, http.StatusPaymentRequired)
		return
	}
	// withdrawn record save and add new order record
	blnc.Outcome = b.Sum
	blnc.OrderNumber = b.Order

	ordr := methods.NewOrder(b.Order)

	if err = blnc.Add(); err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw add new balance record err:%v", err)
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	if err = ordr.Add(); err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw add new order err:%v", err)
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	resp.NoContent(w, http.StatusOK)
}
