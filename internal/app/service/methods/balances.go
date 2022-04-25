package methods

import (
	"database/sql"
	"errors"
	db "github.com/EestiChameleon/GOphermart/internal/app/storage"
	"github.com/EestiChameleon/GOphermart/internal/models"
	"github.com/georgysavva/scany/pgxscan"
	"time"
)

var (
	ErrBalanceInsertFailed = errors.New("failed to save new balance record")
)

type Balance struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProcessedAt time.Time `json:"processed_at"`
	Income      int       `json:"income"`
	Outcome     int       `json:"outcome"`
	OrderNumber string    `json:"order_number"`
}

func NewBalanceRecord(uID int, ordNumber string) *Balance {
	return &Balance{
		ID:          0,
		UserID:      uID,
		ProcessedAt: time.Now(),
		OrderNumber: ordNumber,
	}
}

func (b *Balance) Add() error {
	err := db.Pool.DB.QueryRow(ctx,
		"INSERT INTO balances(user_id, processed_at, income, outcome, order_number) "+
			"VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		b.UserID, b.ProcessedAt, b.Income, b.Outcome, b.OrderNumber).Scan(&b.ID)

	if err != nil {
		return err
	}

	if b.ID < 1 {
		return ErrBalanceInsertFailed
	}

	return nil
}

func GetBalanceAndWithdrawnByUserID(uID int) (*models.BalanceData, error) {
	var c, w sql.NullInt64
	if err := db.Pool.DB.QueryRow(ctx,
		"SELECT sum(income)-sum(outcome) as current, sum(outcome) as withdraw FROM balances WHERE user_id=$1;",
		uID).Scan(&c, &w); err != nil {
		return nil, err
	}

	return &models.BalanceData{
		Current:   float64(c.Int64) / 100,
		Withdrawn: float64(w.Int64) / 100,
	}, nil
}

func GetUserWithdrawals(uID int) ([]*models.WithdrawalsData, error) {
	var list []*models.WithdrawalsData
	err := pgxscan.Select(ctx, db.Pool.DB, &list,
		"SELECT order_number as order, outcome as sum, processed_at "+
			"FROM balances WHERE outcome != 0 AND user_id=$1;", uID)
	if err != nil {
		return nil, err
	}

	return list, nil
}
