package router

import (
	"github.com/EestiChameleon/GOphermart/internal/app/cfg"
	h "github.com/EestiChameleon/GOphermart/internal/app/router/handlers"
	"github.com/EestiChameleon/GOphermart/internal/app/router/mw"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func Start() error {
	// Chi instance
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	//router.Use(mw.GZIP)

	// Routes

	router.Post("/api/user/register", h.UserRegister)
	router.Post("/api/user/login", h.UserLogin)

	router.With(mw.AuthCheck).Route("/api/user", func(r chi.Router) {
		r.Get("/orders", h.UserOrdersList)
		r.Get("/balance", h.UserBalance)
		r.Get("/balance/withdrawals", h.UserBalanceWithdrawals)

		r.Post("/orders", h.UserAddOrder)
		r.Post("/balance/withdraw", h.UserBalanceWithdraw)

	})

	// Start server
	s := http.Server{
		Addr:    cfg.Envs.RunAddr,
		Handler: router,
		// ReadTimeout: 30 * time.Second, // customize http.Server timeouts
	}

	cmlogger.Sug.Infow("SERVER STARTED", "Time", time.Now().Format(time.RFC3339))
	return s.ListenAndServe()
}
