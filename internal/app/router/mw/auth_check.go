package mw

import (
	resp "github.com/EestiChameleon/GOphermart/internal/app/router/responses"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/ctxfunc"
	"net/http"
)

func AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("gophermartID")
		if err != nil {
			resp.NoContent(w, http.StatusUnauthorized)
			return
		}

		userID, err := service.JWTDecodeUserID(cookie.Value)
		if err != nil {
			resp.NoContent(w, http.StatusInternalServerError)
			return
		}

		newCtx := ctxfunc.SetUserIDToCTX(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
