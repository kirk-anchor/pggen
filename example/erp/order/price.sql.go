// Code generated by pggen. DO NOT EDIT.

package order

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

const findOrdersByPriceSQL = `SELECT * FROM orders WHERE order_total > $1;`

type FindOrdersByPriceRow struct {
	OrderID    int32           `json:"order_id"`
	OrderDate  time.Time       `json:"order_date"`
	OrderTotal decimal.Decimal `json:"order_total"`
	CustomerID *int32          `json:"customer_id"`
}

// FindOrdersByPrice implements Querier.FindOrdersByPrice.
func (q *DBQuerier) FindOrdersByPrice(ctx context.Context, minTotal decimal.Decimal) ([]FindOrdersByPriceRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOrdersByPrice")
	rows, err := q.conn.Query(ctx, findOrdersByPriceSQL, minTotal)
	if err != nil {
		return nil, fmt.Errorf("query FindOrdersByPrice: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[FindOrdersByPriceRow])
	return res, q.errWrap(err)
}

const findOrdersMRRSQL = `SELECT date_trunc('month', order_date) AS month, sum(order_total) AS order_mrr
FROM orders
GROUP BY date_trunc('month', order_date);`

type FindOrdersMRRRow struct {
	Month    *time.Time          `json:"month"`
	OrderMRR decimal.NullDecimal `json:"order_mrr"`
}

// FindOrdersMRR implements Querier.FindOrdersMRR.
func (q *DBQuerier) FindOrdersMRR(ctx context.Context) ([]FindOrdersMRRRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOrdersMRR")
	rows, err := q.conn.Query(ctx, findOrdersMRRSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindOrdersMRR: %w", q.errWrap(err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[FindOrdersMRRRow])
	return res, q.errWrap(err)
}
