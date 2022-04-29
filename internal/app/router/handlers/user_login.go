package handlers

import (
	"encoding/json"
	"errors"
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	db "github.com/EestiChameleon/GOphermart/internal/app/storage"
	"github.com/EestiChameleon/GOphermart/internal/models"
	"io"
	"net/http"
)

// UserLogin аутентификация пользователя;
/*
POST /api/user/login HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}
Возможные коды ответа:

200 — пользователь успешно аутентифицирован;
400 — неверный формат запроса;
401 — неверная пара логин/пароль;
500 — внутренняя ошибка сервера.
*/

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var b models.LoginData
	data, err := io.ReadAll(r.Body)
	if err != nil {
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &b)
	if err != nil {
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	if b.Password == "" || b.Login == "" {
		resp.NoContent(w, http.StatusBadRequest) // 401?
		return
	}

	token, err := service.CheckAuthData(b)
	if err != nil && !errors.Is(err, service.ErrWrongAuthData) && !errors.Is(err, db.ErrNotFound) {
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, service.ErrWrongAuthData) || errors.Is(err, db.ErrNotFound) {
		resp.NoContent(w, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, resp.CreateCookie("gophermartID", token))
	resp.NoContent(w, http.StatusOK)
}
