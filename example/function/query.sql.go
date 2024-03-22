// Code generated by pggen. DO NOT EDIT.

package function

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	OutParams(ctx context.Context) ([]OutParamsRow, error)
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
		return nil, fmt.Errorf("query OutParams: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*[]ListItem)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*ListStats)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (OutParamsRow, error) {
		vals := row.RawValues()
		var item OutParamsRow
		if err := plan0.Scan(vals[0], &item.Items); err != nil {
			return item, fmt.Errorf("scan OutParams._items: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan OutParams._stats: %w", err)
		}
		return item, nil
	})
}

type scanCacheKey struct {
	oid      uint32
	format   int16
	typeName string
}

var (
	plans   = make(map[scanCacheKey]pgtype.ScanPlan, 16)
	plansMu sync.RWMutex
)

func planScan(codec pgtype.Codec, fd pgconn.FieldDescription, target any) pgtype.ScanPlan {
	key := scanCacheKey{fd.DataTypeOID, fd.Format, fmt.Sprintf("%T", target)}
	plansMu.RLock()
	plan := plans[key]
	plansMu.RUnlock()
	if plan != nil {
		return plan
	}
	plan = codec.PlanScan(nil, fd.DataTypeOID, fd.Format, target)
	plansMu.Lock()
	plans[key] = plan
	plansMu.Unlock()
	return plan
}

type ptrScanner[T any] struct {
	basePlan pgtype.ScanPlan
}

func (s ptrScanner[T]) Scan(src []byte, dst any) error {
	if src == nil {
		return nil
	}
	d := dst.(**T)
	*d = new(T)
	return s.basePlan.Scan(src, *d)
}

func planPtrScan[T any](codec pgtype.Codec, fd pgconn.FieldDescription, target *T) pgtype.ScanPlan {
	key := scanCacheKey{fd.DataTypeOID, fd.Format, fmt.Sprintf("*%T", target)}
	plansMu.RLock()
	plan := plans[key]
	plansMu.RUnlock()
	if plan != nil {
		return plan
	}
	basePlan := planScan(codec, fd, target)
	ptrPlan := ptrScanner[T]{basePlan}
	plansMu.Lock()
	plans[key] = plan
	plansMu.Unlock()
	return ptrPlan
}
