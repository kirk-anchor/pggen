// Code generated by pggen. DO NOT EDIT.

package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const findOrdersByPriceSQL = `SELECT * FROM orders WHERE order_total > $1;`

type FindOrdersByPriceRow struct {
	OrderID    int32              `json:"order_id"`
	OrderDate  pgtype.Timestamptz `json:"order_date"`
	OrderTotal pgtype.Numeric     `json:"order_total"`
	CustomerID *int32             `json:"customer_id"`
}

// FindOrdersByPrice implements Querier.FindOrdersByPrice.
func (q *DBQuerier) FindOrdersByPrice(ctx context.Context, minTotal pgtype.Numeric) ([]FindOrdersByPriceRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOrdersByPrice")
	rows, err := q.conn.Query(ctx, findOrdersByPriceSQL, minTotal)
	if err != nil {
		return nil, q.errWrap(fmt.Errorf("query FindOrdersByPrice: %w", err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[FindOrdersByPriceRow])
	return res, q.errWrap(err)
}

const findOrdersMRRSQL = `SELECT date_trunc('month', order_date) AS month, sum(order_total) AS order_mrr
FROM orders
GROUP BY date_trunc('month', order_date);`

type FindOrdersMRRRow struct {
	Month    pgtype.Timestamptz `json:"month"`
	OrderMRR pgtype.Numeric     `json:"order_mrr"`
}

// FindOrdersMRR implements Querier.FindOrdersMRR.
func (q *DBQuerier) FindOrdersMRR(ctx context.Context) ([]FindOrdersMRRRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOrdersMRR")
	rows, err := q.conn.Query(ctx, findOrdersMRRSQL)
	if err != nil {
		return nil, q.errWrap(fmt.Errorf("query FindOrdersMRR: %w", err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[FindOrdersMRRRow])
	return res, q.errWrap(err)
}
