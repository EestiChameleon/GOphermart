package methods

import (
	"errors"
	db "github.com/EestiChameleon/GOphermart/internal/app/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"time"
)

var (
	ErrOrderInsertFailed  = errors.New("failed to save new order")
	ErrOrderUpdateFailed  = errors.New("failed to update order")
	ErrOrderWrongOwner    = errors.New("order with provided number owned by another user")
	ErrOrderAlreadyExists = errors.New("order with provided number already exists in database")
)

type Order struct {
	Number     string    `json:"number"`
	UserID     int       `json:"user_id"`
	UploadedAt time.Time `json:"uploaded_at"` // my time type
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
}

func NewOrder(uID int, number string) *Order {
	return &Order{
		Number:     number,
		UserID:     uID,
		UploadedAt: time.Now(),
		Status:     "NEW",
	}
}

// CheckNumber verify the order number. 1) Check for existence. 2) Check for correct UserID
//
func (o *Order) CheckNumber() error {
	var dbUserID int
	err := db.Pool.DB.QueryRow(ctx, "SELECT user_id FROM orders WHERE number=$1", o.Number).Scan(&dbUserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		// order not found - ok
		return nil
	}

	// order found - check userID match
	if dbUserID != o.UserID {
		return ErrOrderWrongOwner
	}

	return ErrOrderAlreadyExists
}

func (o *Order) GetByNumber() error {
	err := pgxscan.Get(ctx, db.Pool.DB, o,
		"SELECT user_id, uploaded_at, status FROM orders WHERE number=$1", o.Number)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return db.ErrNotFound
	}

	return nil
}

func (o *Order) Add() error {
	tag, err := db.Pool.DB.Exec(ctx,
		"INSERT INTO orders(number, user_id, uploaded_at, status) "+
			"VALUES ($1, $2, $3, $4) ON CONFLICT (number) DO NOTHING;",
		o.Number, o.UserID, o.UploadedAt, o.Status)

	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrOrderInsertFailed
	}

	return nil
}

func (o *Order) UpdateStatus(status string) error {
	tag, err := db.Pool.DB.Exec(ctx,
		"UPDATE orders SET status = $1 WHERE number = $2;",
		status, o.Number)

	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrOrderUpdateFailed
	}

	return nil
}

func (o *Order) SetAccrual(value int) error {
	tag, err := db.Pool.DB.Exec(ctx,
		"UPDATE orders SET accrual = $2 WHERE number = $1;",
		o.Number, value)

	if err != nil {
		return err
	}

	if tag.RowsAffected() < 1 {
		return ErrOrderUpdateFailed
	}

	return nil
}

func (o *Order) UpdStatusSetAccrual(status string, value int) error {
	if err := o.UpdateStatus(status); err != nil {
		return err
	}

	if err := o.SetAccrual(value); err != nil {
		return err
	}

	return nil
}

func GetOrdersListByUserID(uID int) ([]*Order, error) {
	var list []*Order
	err := pgxscan.Select(ctx, db.Pool.DB, &list,
		"SELECT * FROM orders WHERE user_id=$1", uID)
	return list, err
}

func GetOrdersListNotFinal() ([]*Order, error) {
	var list []*Order
	err := pgxscan.Select(ctx, db.Pool.DB, &list,
		"SELECT * FROM orders WHERE status in ('NEW', 'PROCESSING');")
	return list, err
}
