package service

import (
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"log"
	"net/http"
	"time"
)

/*
Доступные статусы обработки заказов:
NEW — заказ загружен в систему, но не попал в обработку. Статус проставляется при первичном попадании в БД.
PROCESSING — вознаграждение за заказ рассчитывается. Статус проставляется при получении статусов REGISTERED & PROCESSING
INVALID — система расчёта вознаграждений отказала в расчёте;
PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
*/

func PollOrderCron(accrualClient AccrualSystem, cronPeriod time.Duration) {
	ticker := time.NewTicker(cronPeriod)

	for {
		select {
		case <-ticker.C:
			if err := proccessOrders(accrualClient); err != nil {
				log.Println("PollOrderCron err:", err)
				continue
			}
		}
	}
}

func proccessOrders(accrualClient AccrualSystem) error {
	// get all orders with NOT final status
	orders, err := methods.GetOrdersListNotFinal()
	if err != nil {
		log.Println("select orders from DB err:", err)
		return err
	}

	for _, order := range orders {
		switch order.Status {
		case "NEW":
			o := methods.NewOrder(order.Number)
			if err = o.UpdateStatus("PROCESSING"); err != nil {
				log.Printf("update order #%s failed: %v", order.Number, err)
			}
			continue

		case "PROCESSING":
			orderInfo := accrualClient.GetOrderInfo(order.Number)
			if accrualClient.ReturnStatus() == http.StatusOK {
				// successful request
				if orderInfo.Status == "INVALID" {
					o := methods.NewOrder(order.Number)
					if err = o.UpdateStatus("INVALID"); err != nil {
						log.Printf("update order #%s failed: %v", order.Number, err)
						continue
					}
					b := methods.NewBalanceRecord()
					b.OrderNumber = o.Number
					b.Income = o.Accrual
				}
				if orderInfo.Status == "PROCESSED" {
					o := methods.NewOrder(order.Number)
					if err = o.SetProcessedAndAccrual(order.Accrual); err != nil {
						log.Printf("update order #%s failed: %v", order.Number, err)
						continue
					}
				}
			} else {
				log.Println("accrual system response status code:", accrualClient.ReturnStatus())
			}
		}
	}

	return nil
}
