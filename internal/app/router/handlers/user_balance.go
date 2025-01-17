package handlers

import (
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/EestiChameleon/GOphermart/internal/ctxfunc"
	"net/http"
)

// UserBalance предоставляет возможность получения текущего баланса счёта баллов лояльности пользователя
/*
GET /api/user/balance HTTP/1.1
Content-Length: 0
Возможные коды ответа:
200 — успешная обработка запроса.
Формат ответа:
200 OK HTTP/1.1
Content-Type: application/json
{
	"current": 500.5,
	"withdrawn": 42
}
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func UserBalance(w http.ResponseWriter, r *http.Request) {
	userID := ctxfunc.GetUserIDFromCTX(r.Context())
	if userID < 1 {
		resp.NoContent(w, http.StatusUnauthorized)
		return
	}

	res, err := methods.GetBalanceAndWithdrawnByUserID(userID)
	if err != nil {
		cmlogger.Sug.Errorf("user balance err:%v", err)
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	cmlogger.Sug.Infow("get user balance", "UserID", userID, "balance", res)
	resp.JSON(w, http.StatusOK, res)
}
