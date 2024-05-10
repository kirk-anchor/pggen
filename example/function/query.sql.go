// Code generated by pggen. DO NOT EDIT.

package function

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	OutParams(ctx context.Context) ([]OutParamsRow, error)
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

// ListItem represents the Postgres composite type "list_item".
type ListItem struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

// ListStats represents the Postgres composite type "list_stats".
type ListStats struct {
	Val1 *string  `json:"val1"`
	Val2 []*int32 `json:"val2"`
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

var _ = addTypeToRegister("list_item")

var _ = addTypeToRegister("list_stats")

var _ = addTypeToRegister("_list_item")

const outParamsSQL = `SELECT * FROM out_params();`

type OutParamsRow struct {
	Items []ListItem `json:"_items"`
	Stats ListStats  `json:"_stats"`
}

// OutParams implements Querier.OutParams.
func (q *DBQuerier) OutParams(ctx context.Context) ([]OutParamsRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "OutParams")
	rows, err := q.conn.Query(ctx, outParamsSQL)
	if err != nil {
		return nil, q.errWrap(fmt.Errorf("query OutParams: %w", err))
	}
	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[OutParamsRow])
	return res, q.errWrap(err)
}
