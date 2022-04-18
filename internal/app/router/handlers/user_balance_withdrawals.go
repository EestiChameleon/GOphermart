package handlers

import (
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/models"
	"net/http"
)

// UserBalanceWithdrawals получение информации о выводе средств с накопительного счёта пользователем
/*
Формат запроса:
GET /api/user/withdrawals HTTP/1.1
Content-Length: 0
Возможные коды ответа:
200 — успешная обработка запроса.
Формат ответа:
  200 OK HTTP/1.1
  Content-Type: application/json
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
func UserBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {
	var ubw []*models.WithdrawalsData
	if err := methods.GetUserWithdrawals(&ubw); err != nil {
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	if len(ubw) < 1 {
		resp.NoContent(w, http.StatusNoContent)
		return
	}

	resp.JSON(w, http.StatusOK, ubw)
}
