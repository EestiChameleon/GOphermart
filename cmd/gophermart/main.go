package main

import (
	"github.com/EestiChameleon/GOphermart/internal/app/cfg"
	"github.com/EestiChameleon/GOphermart/internal/app/router"
	"github.com/EestiChameleon/GOphermart/internal/app/service"
	"github.com/EestiChameleon/GOphermart/internal/app/storage"
	"log"
	"time"
)

func main() {
	// parsing of the environments + flags
	if err := cfg.GetEnvs(); err != nil {
		log.Fatal("Env parse err:", err)
	}
	log.Println("envs parsed")

	// init DB connection + migrations
	if err := storage.InitConnection(); err != nil {
		log.Fatal("Storage init err:", err)
	}
	defer storage.Shutdown()
	log.Println("DB connected")

	// init accrual instance
	service.AccrualBot = service.NewAccrualClient(cfg.Envs.AccrualSysAddr)
	log.Println("accrual bot initiated. Address:", cfg.Envs.AccrualSysAddr)

	// start the order check loop
	go service.PollOrderCron(service.AccrualBot, time.Second*60)
	log.Println("PollOrderCron launched")

	// start the service
	if err := router.Start(); err != nil {
		log.Fatal("server init err:", err)
	}
}
