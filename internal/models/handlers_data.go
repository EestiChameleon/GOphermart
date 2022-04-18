package models

import (
	"github.com/shopspring/decimal"
	"time"
)

// LoginData - структура данных логин/пароль пользователя
type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// BalanceData - структура для запроса на вывод данных о текущем состоянии бонусного счета пользователя
type BalanceData struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}

// WithdrawData - структура входящего запроса на списание бонусов в счет оплаты заказа
type WithdrawData struct {
	Order string              `json:"order"`
	Sum   decimal.NullDecimal `json:"sum"`
}

// WithdrawalsData - структура для запроса на вывод данных обо всех операциях списания бонусов
type WithdrawalsData struct {
	Order       string              `json:"order"`
	Sum         decimal.NullDecimal `json:"sum"`
	ProcessedAt time.Time           `json:"processed_at"`
}
