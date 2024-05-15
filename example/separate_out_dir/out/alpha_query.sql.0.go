// Code generated by pggen. DO NOT EDIT.

package out

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	AlphaNested(ctx context.Context) (string, error)

	AlphaCompositeArray(ctx context.Context) ([]Alpha, error)

	Alpha(ctx context.Context) (string, error)

	Bravo(ctx context.Context) (string, error)
}

var _ Querier = &DBQuerier{}

type DBQuerier struct {
	conn    genericConn
	errWrap func(err error) error
}

// genericConn is a connection like *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
		errWrap: func(err error) error {
			return err
		},
	}
}

// Alpha represents the Postgres composite type "alpha".
type Alpha struct {
	Key *string `json:"key"`
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

var _ = addTypeToRegister("public.alpha")

var _ = addTypeToRegister("public._alpha")

const alphaNestedSQL = `SELECT 'alpha_nested' as output;`

// AlphaNested implements Querier.AlphaNested.
func (q *DBQuerier) AlphaNested(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "AlphaNested")
	rows, err := q.conn.Query(ctx, alphaNestedSQL)
	if err != nil {
		return "", fmt.Errorf("query AlphaNested: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	return res, q.errWrap(err)
}

const alphaCompositeArraySQL = `SELECT ARRAY[ROW('key')]::alpha[];`

// AlphaCompositeArray implements Querier.AlphaCompositeArray.
func (q *DBQuerier) AlphaCompositeArray(ctx context.Context) ([]Alpha, error) {
	ctx = context.WithValue(ctx, QueryName{}, "AlphaCompositeArray")
	rows, err := q.conn.Query(ctx, alphaCompositeArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query AlphaCompositeArray: %w", q.errWrap(err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[[]Alpha])
	return res, q.errWrap(err)
}
