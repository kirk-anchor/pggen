// Code generated by pggen. DO NOT EDIT.

package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	CreateTenant(ctx context.Context, key string, name string) (CreateTenantRow, error)

	FindOrdersByCustomer(ctx context.Context, customerID int32) ([]FindOrdersByCustomerRow, error)

	FindProductsInOrder(ctx context.Context, orderID int32) ([]FindProductsInOrderRow, error)

	InsertCustomer(ctx context.Context, params InsertCustomerParams) (InsertCustomerRow, error)

	InsertOrder(ctx context.Context, params InsertOrderParams) (InsertOrderRow, error)

	FindOrdersByPrice(ctx context.Context, minTotal pgtype.Numeric) ([]FindOrdersByPriceRow, error)

	FindOrdersMRR(ctx context.Context) ([]FindOrdersMRRRow, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn genericConn
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{conn: conn}
}

// RegisterTypes should be run in config.AfterConnect to load custom types
func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	for _, typ := range typesToRegister {
		dt, err := conn.LoadType(ctx, typ)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(dt)
	}
	return nil
}

var typesToRegister = []string{}

func addTypeToRegister(typ string) struct{} {
	typesToRegister = append(typesToRegister, typ)
	return struct{}{}
}

const createTenantSQL = `INSERT INTO tenant (tenant_id, name)
VALUES (base36_decode($1::text)::tenant_id, $2::text)
RETURNING *;`

type CreateTenantRow struct {
	TenantID int     `json:"tenant_id"`
	Rname    *string `json:"rname"`
	Name     string  `json:"name"`
}

// CreateTenant implements Querier.CreateTenant.
func (q *DBQuerier) CreateTenant(ctx context.Context, key string, name string) (CreateTenantRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CreateTenant")
	rows, err := q.conn.Query(ctx, createTenantSQL, key, name)
	if err != nil {
		return CreateTenantRow{}, fmt.Errorf("query CreateTenant: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CreateTenantRow])
}

const findOrdersByCustomerSQL = `SELECT *
FROM orders
WHERE customer_id = $1;`

type FindOrdersByCustomerRow struct {
	OrderID    int32              `json:"order_id"`
	OrderDate  pgtype.Timestamptz `json:"order_date"`
	OrderTotal pgtype.Numeric     `json:"order_total"`
	CustomerID *int32             `json:"customer_id"`
}

// FindOrdersByCustomer implements Querier.FindOrdersByCustomer.
func (q *DBQuerier) FindOrdersByCustomer(ctx context.Context, customerID int32) ([]FindOrdersByCustomerRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOrdersByCustomer")
	rows, err := q.conn.Query(ctx, findOrdersByCustomerSQL, customerID)
	if err != nil {
		return nil, fmt.Errorf("query FindOrdersByCustomer: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[FindOrdersByCustomerRow])
}

const findProductsInOrderSQL = `SELECT o.order_id, p.product_id, p.name
FROM orders o
  INNER JOIN order_product op USING (order_id)
  INNER JOIN product p USING (product_id)
WHERE o.order_id = $1;`

type FindProductsInOrderRow struct {
	OrderID   *int32  `json:"order_id"`
	ProductID *int32  `json:"product_id"`
	Name      *string `json:"name"`
}

// FindProductsInOrder implements Querier.FindProductsInOrder.
func (q *DBQuerier) FindProductsInOrder(ctx context.Context, orderID int32) ([]FindProductsInOrderRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindProductsInOrder")
	rows, err := q.conn.Query(ctx, findProductsInOrderSQL, orderID)
	if err != nil {
		return nil, fmt.Errorf("query FindProductsInOrder: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[FindProductsInOrderRow])
}

const insertCustomerSQL = `INSERT INTO customer (first_name, last_name, email)
VALUES ($1, $2, $3)
RETURNING *;`

type InsertCustomerParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type InsertCustomerRow struct {
	CustomerID int32  `json:"customer_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
}

// InsertCustomer implements Querier.InsertCustomer.
func (q *DBQuerier) InsertCustomer(ctx context.Context, params InsertCustomerParams) (InsertCustomerRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertCustomer")
	rows, err := q.conn.Query(ctx, insertCustomerSQL, params.FirstName, params.LastName, params.Email)
	if err != nil {
		return InsertCustomerRow{}, fmt.Errorf("query InsertCustomer: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[InsertCustomerRow])
}

const insertOrderSQL = `INSERT INTO orders (order_date, order_total, customer_id)
VALUES ($1, $2, $3)
RETURNING *;`

type InsertOrderParams struct {
	OrderDate  pgtype.Timestamptz `json:"order_date"`
	OrderTotal pgtype.Numeric     `json:"order_total"`
	CustID     int32              `json:"cust_id"`
}

type InsertOrderRow struct {
	OrderID    int32              `json:"order_id"`
	OrderDate  pgtype.Timestamptz `json:"order_date"`
	OrderTotal pgtype.Numeric     `json:"order_total"`
	CustomerID *int32             `json:"customer_id"`
}

// InsertOrder implements Querier.InsertOrder.
func (q *DBQuerier) InsertOrder(ctx context.Context, params InsertOrderParams) (InsertOrderRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertOrder")
	rows, err := q.conn.Query(ctx, insertOrderSQL, params.OrderDate, params.OrderTotal, params.CustID)
	if err != nil {
		return InsertOrderRow{}, fmt.Errorf("query InsertOrder: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[InsertOrderRow])
}
