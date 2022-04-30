package models

import (
	"time"
)

// LoginData - структура данных логин/пароль пользователя
type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// BalanceData - структура для запроса на вывод данных о текущем состоянии бонусного счета пользователя
type BalanceData struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// WithdrawData - структура входящего запроса на списание бонусов в счет оплаты заказа
type WithdrawData struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// WithdrawalsData - структура для запроса на вывод данных обо всех операциях списания бонусов
type WithdrawalsData struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
