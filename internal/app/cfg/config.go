package cfg

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS" envDefault:"localhost:8080"` // Адрес запуска HTTP-сервера
	DatabaseURI    string `env:"DATABASE_URI"`                            // Строка с адресом подключения к БД
	AccrualSysAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`                  // Адрес системы расчёта начислений «http://server:port»

	CryptoKey    string `env:"CRYPTO_KEY"`    // secret word to encrypt/decrypt JWT for cookies envDefault:"secret_123456789"
	CronCooldown int    `env:"CRON_COOLDOWN"` //
}

var Envs Config

type ContextKey string

func GetEnvs() error {
	flag.StringVar(&Envs.RunAddr, "a", "http://localhost:8080", "RUN_ADDRESS to listen on")
	flag.StringVar(&Envs.DatabaseURI, "d", "", "DATABASE_DSN. Address for connection to DB")
	flag.StringVar(&Envs.AccrualSysAddr, "r", "", "ACCRUAL_SYSTEM_ADDRESS ")

	flag.IntVar(&Envs.CronCooldown, "cd", 1, "CRON_COOLDOWN")

	if err := env.Parse(&Envs); err != nil {
		return err
	}

	flag.Parse()
	log.Printf("vale %v type %T", Envs.CronCooldown, Envs.CronCooldown)
	return nil
}
