package handlers

import (
	"encoding/json"
	"errors"
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	s "github.com/EestiChameleon/GOphermart/internal/app/service"
	m "github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/models"
	"io"
	"net/http"
)

// UserRegister регистрация пользователя
/*
POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
	"login": "<login>",
	"password": "<password>"
}

Возможные коды ответа:

200 — пользователь успешно зарегистрирован и аутентифицирован;
400 — неверный формат запроса;
409 — логин уже занят;
500 — внутренняя ошибка сервера.
*/

func UserRegister(w http.ResponseWriter, r *http.Request) {
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
		resp.NoContent(w, http.StatusBadRequest)
		return
	}

	u := m.NewUser(b.Login, s.EncryptPass(b.Password))

	if err = u.Add(); err != nil {
		if errors.Is(err, m.ErrLoginUnavailable) {
			resp.NoContent(w, http.StatusConflict)
			return
		}
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	token, err := s.JWTEncodeUserID(u.ID)
	if err != nil {
		resp.NoContent(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, resp.CreateCookie("gophermartID", token))
	resp.NoContent(w, http.StatusOK)
}
