package service

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"net/http"
)

// GetOrderAccrualInfo is a method, that makes a request GET /api/orders/{orderNumber} to the Accrual System
// As response receives json
// {
//      "order": "<number>",
//      "status": "PROCESSED",
//      "accrual": 500
//  }
/*
Возможные коды ответа:
200 — успешная обработка запроса.
Формат ответа:
  200 OK HTTP/1.1
  Content-Type: application/json
  {
      "order": "<number>",		order — номер заказа
      "status": "PROCESSED", 	status — статус расчёта начисления
      "accrual": 500			accrual — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе
  }

Status:
REGISTERED — заказ зарегистрирован, но не начисление не рассчитано;
INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
PROCESSING — расчёт начисления в процессе;
PROCESSED — расчёт начисления окончен;

Errors:
429 — превышено количество запросов к сервису.
Формат ответа:
  429 Too Many Requests HTTP/1.1
  Content-Type: text/plain
  Retry-After: 60

  No more than N requests per minute allowed

500 — внутренняя ошибка сервера.
Заказ может быть взят в расчёт в любой момент после его совершения. Время выполнения расчёта системой не регламентировано
Статусы INVALID и PROCESSED являются окончательными.
Общее количество запросов информации о начислении не ограничено.
*/

var (
	AccrualBot AccrualSystem
)

type AccrualSystem interface {
	GetOrderInfo(orderNumber string) *OrderAccrualInfo
	ReturnStatus() int
}

// OrderAccrualInfo - response
type OrderAccrualInfo struct {
	Order   string          `json:"order"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}

type AccrualClient struct {
	AccrualSystemAddress string
	RespStatusCode       int
}

func NewAccrualClient(address string) *AccrualClient {
	return &AccrualClient{
		AccrualSystemAddress: address,
	}
}

func (ac *AccrualClient) GetOrderInfo(orderNumber string) *OrderAccrualInfo {
	client := http.Client{}
	accSysPath := ac.AccrualSystemAddress + "/api/orders/" + orderNumber
	getReq, err := http.NewRequest(http.MethodGet, accSysPath, nil)
	if err != nil {
		log.Println("Accrual System NEW GET request err:", err)
	}
	res, err := client.Do(getReq)
	if err != nil {
		log.Println("Accrual System DO GET request err:", err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusTooManyRequests:
		ac.RespStatusCode = http.StatusTooManyRequests
		return nil
	case http.StatusInternalServerError:
		ac.RespStatusCode = http.StatusInternalServerError
		return nil
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		ac.RespStatusCode = http.StatusInternalServerError
		return nil
	}

	orderInfo := new(OrderAccrualInfo)

	err = json.Unmarshal(data, &orderInfo)
	if err != nil {
		ac.RespStatusCode = http.StatusInternalServerError
		return nil
	}

	ac.RespStatusCode = res.StatusCode
	return orderInfo
}

func (ac *AccrualClient) ReturnStatus() int {
	return ac.RespStatusCode
}
