package service

import (
	"errors"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/cmlogger"
	"github.com/EestiChameleon/GOphermart/internal/models"
)

var (
	ErrWithdrawUnavailable = errors.New("current bonus balance is lower than the required withdraw sum")
)

// BalanceWithdraw verify the withdraw possibility and in case of success save new order + new outcome balance record
func BalanceWithdraw(uID int, wd models.WithdrawData) error {
	// get current balance and whole withdrawn
	res, err := methods.GetBalanceAndWithdrawnByUserID(uID)
	if err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw GetBalanceAndWithdrawnByUserID err:%v", err)
		return err
	}

	// withdraw possibility check
	if res.Current < wd.Sum {
		cmlogger.Sug.Infow("current bonus balance is lower than required bonus",
			"current", res.Current, "bonus required", res.Withdrawn, "status", "REFUSED")
		return ErrWithdrawUnavailable
	}

	// withdrawn sum save and add new order record. Convert sum float to int
	balance := methods.NewBalanceRecord(uID, wd.Order)
	balance.Outcome = int(wd.Sum * 100) // 758.99 -> 75899
	if err = balance.Add(); err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw add new balance record err:%v", err)
		return err
	}

	order := methods.NewOrder(uID, wd.Order)
	if err = order.Add(); err != nil {
		cmlogger.Sug.Errorf("UserBalanceWithdraw add new order err:%v", err)
		return err
	}

	return nil
}
