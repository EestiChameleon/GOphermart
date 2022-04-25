package accrual

import (
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"math/rand"
	"time"
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

const (
	TestOrderStatusRegistered = "REGISTERED"
	TestOrderStatusInvalid    = "INVALID"
	TestOrderStatusProcessing = "PROCESSING"
	TestOrderStatusProcessed  = "PROCESSED"
)

type TestAccrualClient struct {
	AccrualSystemAddress string
}

func NewTestAccrualClient(address string) *TestAccrualClient {
	return &TestAccrualClient{
		AccrualSystemAddress: address,
	}
}

func (ac *TestAccrualClient) GetOrderInfo(orderNumber string) (*OrderAccrualInfo, error) {
	x := GetRand(2)
	switch x {
	case 3: // processing / registered
		cmlogger.Sug.Infow("order info case", "Number", orderNumber, "Status", TestOrderStatusProcessing)
		return newOrderInfo(orderNumber, TestOrderStatusProcessing), nil
	case 2: // invalid
		cmlogger.Sug.Infow("order info case", "Number", orderNumber, "Status", TestOrderStatusInvalid)
		return newOrderInfo(orderNumber, TestOrderStatusInvalid), nil
	case 1: // processed
		order := newOrderInfo(orderNumber, TestOrderStatusProcessed)
		order.Accrual = GetAccrual()
		cmlogger.Sug.Infow("order info case", "Number", orderNumber, "Status", TestOrderStatusProcessed, "Accrual", order.Accrual)
		return order, nil
	case 0:
		return nil, ErrAccSysTooManyReq
	default:
		cmlogger.Sug.Infow("order info case", "Number", orderNumber, "Status", TestOrderStatusRegistered)
		return newOrderInfo(orderNumber, TestOrderStatusRegistered), nil
	}
}

func newOrderInfo(number, status string) *OrderAccrualInfo {
	return &OrderAccrualInfo{
		Order:  number,
		Status: status,
	}
}

func GetRand(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func GetAccrual() float64 {
	rand.Seed(time.Now().UnixNano())

	return float64(rand.Intn(100000)) / 100
}
