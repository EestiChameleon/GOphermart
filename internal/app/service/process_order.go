package service

import (
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/EestiChameleon/GOphermart/internal/pkg/accrual"
	"log"
	"time"
)

/*
Доступные статусы обработки заказов:
NEW — заказ загружен в систему, но не попал в обработку. Статус проставляется при первичном попадании в БД.
PROCESSING — вознаграждение за заказ рассчитывается. Статус проставляется при получении статусов REGISTERED & PROCESSING
INVALID — система расчёта вознаграждений отказала в расчёте;
PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
*/

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

func PollOrderCron(accrualClient accrual.AccrualSystem, cronPeriod time.Duration) {
	ticker := time.NewTicker(cronPeriod)
	// instant call
	if err := processOrders(accrualClient); err != nil {
		log.Println("First PollOrderCron call err:", err)
	}
	// then - loop with timer
	for range ticker.C {
		if err := processOrders(accrualClient); err != nil {
			log.Println("Loop PollOrderCron err:", err)
			continue
		}
	}
}

func processOrders(accrualClient accrual.AccrualSystem) error {
	// get all orders with NOT final status
	orders, err := methods.GetOrdersListNotFinal()
	if err != nil {
		log.Println("select orders from DB err:", err)
		return err
	}

	cmlogger.Sug.Infow("--- orders info ---", "Orders", orders)

	for _, order := range orders {
		switch order.Status {
		case OrderStatusNew:
			if err = order.UpdateStatus(OrderStatusProcessing); err != nil {
				log.Printf("update order #%s failed: %v", order.Number, err)
			}
			continue

		case OrderStatusProcessing:
			orderInfo, err := accrualClient.GetOrderInfo(order.Number)
			if err != nil {
				cmlogger.Sug.Infow("accrual system err", "error", err)
				continue
			}

			// successful request
			if orderInfo.Status == accrual.OrderStatusInvalid {
				if err = InvalidOrder(order); err != nil {
					cmlogger.Sug.Infow("update invalid order failed", "order", order.Number, "err", err)
					continue
				}
			}
			if orderInfo.Status == accrual.OrderStatusProcessed {
				if err = ProcessedOrder(order, int(orderInfo.Accrual*100)); err != nil {
					cmlogger.Sug.Infow("update processed order failed",
						"order", order.Number,
						"accrual", order.Accrual,
						"err", err)
				}
			}
		}
	}

	return nil
}

func InvalidOrder(order *methods.Order) error {
	return order.UpdateStatus(OrderStatusInvalid)
}

func ProcessedOrder(order *methods.Order, accrualValue int) error {
	cmlogger.Sug.Infow("PROCESSED order", "Number", order.Number, "Accrual", accrualValue)
	if err := order.UpdStatusSetAccrual(OrderStatusProcessed, accrualValue); err != nil {
		return err
	}
	order.Accrual = accrualValue

	b := methods.NewBalanceRecord(order.UserID, order.Number)
	b.Income = accrualValue
	if err := b.Add(); err != nil {
		return err
	}

	return nil
}
