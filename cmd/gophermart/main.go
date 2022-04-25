package main

import (
	"github.com/EestiChameleon/GOphermart/internal/app/cfg"
	"github.com/EestiChameleon/GOphermart/internal/app/router"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/app/storage"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/EestiChameleon/GOphermart/internal/pkg/accrual"
	"time"
)

func main() {
	cmlogger.InitLogger()

	// parsing of the environments + flags
	if err := cfg.GetEnvs(); err != nil {
		cmlogger.Sug.Fatal("Env parse err:", err)
	}
	cmlogger.Sug.Info("envs parsed")

	// init DB connection + migrations
	if err := storage.InitConnection(); err != nil {
		cmlogger.Sug.Fatal("Storage init err:", err)
	}
	defer storage.Shutdown()
	cmlogger.Sug.Info("DB connected")

	//// init accrual instance
	if cfg.Envs.AccrualSysAddr != "http://localhost:666" {
		accrual.AccrualBot = accrual.NewAccrualClient(cfg.Envs.AccrualSysAddr)
		cmlogger.Sug.Infow("accrual bot initiated", "Address:", cfg.Envs.AccrualSysAddr)
	} else {
		// test accrual instance
		accrual.AccrualBot = accrual.NewTestAccrualClient(cfg.Envs.AccrualSysAddr)
		cmlogger.Sug.Infow("test accrual bot initiated", "Address:", cfg.Envs.AccrualSysAddr)
	}

	// start the order check loop
	go service.PollOrderCron(accrual.AccrualBot, time.Second*1)
	cmlogger.Sug.Info("PollOrderCron launched")

	// start the service
	if err := router.Start(); err != nil {
		cmlogger.Sug.Fatal("server init err:", err)
	}
}
