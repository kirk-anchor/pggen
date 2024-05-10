// Code generated by pggen. DO NOT EDIT.

package out

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const alphaSQL = `SELECT 'alpha' as output;`

// Alpha implements Querier.Alpha.
func (q *DBQuerier) Alpha(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "Alpha")
	rows, err := q.conn.Query(ctx, alphaSQL)
	if err != nil {
		return "", q.errWrap(fmt.Errorf("query Alpha: %w", err))
	}
	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	return res, q.errWrap(err)
}
