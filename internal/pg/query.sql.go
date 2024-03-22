// Code generated by pggen. DO NOT EDIT.

package pg

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
	FindEnumTypes(ctx context.Context, oids []uint32) ([]FindEnumTypesRow, error)

	FindArrayTypes(ctx context.Context, oids []uint32) ([]FindArrayTypesRow, error)

	// A composite type represents a row or record, defined implicitly for each
	// table, or explicitly with CREATE TYPE.
	// https://www.postgresql.org/docs/13/rowtypes.html
	FindCompositeTypes(ctx context.Context, oids []uint32) ([]FindCompositeTypesRow, error)

	// Recursively expands all given OIDs to all descendants through composite
	// types.
	FindDescendantOIDs(ctx context.Context, oids []uint32) ([]uint32, error)

	FindOIDByName(ctx context.Context, name string) (uint32, error)

	FindOIDName(ctx context.Context, oid uint32) (string, error)

	FindOIDNames(ctx context.Context, oid []uint32) ([]FindOIDNamesRow, error)
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

const findEnumTypesSQL = `WITH enums AS (
  SELECT
    enumtypid::int8                                   AS enum_type,
    -- pg_enum row identifier.
    -- The OIDs for pg_enum rows follow a special rule: even-numbered OIDs
    -- are guaranteed to be ordered in the same way as the sort ordering of
    -- their enum type. That is, if two even OIDs belong to the same enum
    -- type, the smaller OID must have the smaller enumsortorder value.
    -- Odd-numbered OID values need bear no relationship to the sort order.
    -- This rule allows the enum comparison routines to avoid catalog
    -- lookups in many common cases. The routines that create and alter enum
    -- types attempt to assign even OIDs to enum values whenever possible.
    array_agg(oid::int8 ORDER BY enumsortorder)       AS enum_oids,
    -- The sort position of this enum value within its enum type. Starts as
    -- 1..n but can be fractional or negative.
    array_agg(enumsortorder ORDER BY enumsortorder)   AS enum_orders,
    -- The textual label for this enum value
    array_agg(enumlabel::text ORDER BY enumsortorder) AS enum_labels
  FROM pg_enum
  GROUP BY pg_enum.enumtypid)
SELECT
  typ.oid           AS oid,
  -- typename: Data type name.
  typ.typname::text AS type_name,
  enum.enum_oids    AS child_oids,
  enum.enum_orders  AS orders,
  enum.enum_labels  AS labels,
  -- typtype: b for a base type, c for a composite type (e.g., a table's
  -- row type), d for a domain, e for an enum type, p for a pseudo-type,
  -- or r for a range type.
  typ.typtype       AS type_kind,
  -- typdefault is null if the type has no associated default value. If
  -- typdefaultbin is not null, typdefault must contain a human-readable
  -- version of the default expression represented by typdefaultbin. If
  -- typdefaultbin is null and typdefault is not, then typdefault is the
  -- external representation of the type's default value, which can be fed
  -- to the type's input converter to produce a constant.
  COALESCE(typ.typdefault, '')    AS default_expr
FROM pg_type typ
  JOIN enums enum ON typ.oid = enum.enum_type
WHERE typ.typisdefined
  AND typ.typtype = 'e'
  AND typ.oid = ANY ($1::oid[]);`

type FindEnumTypesRow struct {
	OID         uint32    `json:"oid"`
	TypeName    string    `json:"type_name"`
	ChildOIDs   []int     `json:"child_oids"`
	Orders      []float32 `json:"orders"`
	Labels      []string  `json:"labels"`
	TypeKind    byte      `json:"type_kind"`
	DefaultExpr string    `json:"default_expr"`
}

// FindEnumTypes implements Querier.FindEnumTypes.
func (q *DBQuerier) FindEnumTypes(ctx context.Context, oids []uint32) ([]FindEnumTypesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindEnumTypes")
	rows, err := q.conn.Query(ctx, findEnumTypesSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindEnumTypes: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*uint32)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*string)(nil))
	plan2 := planScan(pgtype.TextCodec{}, fds[2], (*[]int)(nil))
	plan3 := planScan(pgtype.TextCodec{}, fds[3], (*[]float32)(nil))
	plan4 := planScan(pgtype.TextCodec{}, fds[4], (*[]string)(nil))
	plan5 := planScan(pgtype.TextCodec{}, fds[5], (*byte)(nil))
	plan6 := planScan(pgtype.TextCodec{}, fds[6], (*string)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindEnumTypesRow, error) {
		vals := row.RawValues()
		var item FindEnumTypesRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.oid: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.type_name: %w", err)
		}
		if err := plan2.Scan(vals[2], &item.ChildOIDs); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.child_oids: %w", err)
		}
		if err := plan3.Scan(vals[3], &item.Orders); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.orders: %w", err)
		}
		if err := plan4.Scan(vals[4], &item.Labels); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.labels: %w", err)
		}
		if err := plan5.Scan(vals[5], &item); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.type_kind: %w", err)
		}
		if err := plan6.Scan(vals[6], &item); err != nil {
			return item, fmt.Errorf("scan FindEnumTypes.default_expr: %w", err)
		}
		return item, nil
	})
}

const findArrayTypesSQL = `SELECT
  arr_typ.oid           AS oid,
  -- typename: Data type name.
  arr_typ.typname::text AS type_name,
  elem_typ.oid          AS elem_oid,
  -- typtype: b for a base type, c for a composite type (e.g., a table's
  -- row type), d for a domain, e for an enum type, p for a pseudo-type,
  -- or r for a range type.
  arr_typ.typtype       AS type_kind
FROM pg_type arr_typ
  JOIN pg_type elem_typ ON arr_typ.typelem = elem_typ.oid
WHERE arr_typ.typisdefined
  AND arr_typ.typtype = 'b' -- Array types are base types
  -- If typelem is not 0 then it identifies another row in pg_type. The current
  -- type can then be subscripted like an array yielding values of type typelem.
  -- A “true” array type is variable length (typlen = -1), but some
  -- fixed-length (typlen > 0) types also have nonzero typelem, for example
  -- name and point. If a fixed-length type has a typelem then its internal
  -- representation must be some number of values of the typelem data type with
  -- no other data. Variable-length array types have a header defined by the
  -- array subroutines.
  AND arr_typ.typelem > 0
  -- For a fixed-size type, typlen is the number of bytes in the internal
  -- representation of the type. But for a variable-length type, typlen is
  -- negative. -1 indicates a "varlena" type (one that has a length word), -2
  -- indicates a null-terminated C string.
  AND arr_typ.typlen = -1
  AND arr_typ.oid = ANY ($1::oid[]);`

type FindArrayTypesRow struct {
	OID      uint32 `json:"oid"`
	TypeName string `json:"type_name"`
	ElemOID  uint32 `json:"elem_oid"`
	TypeKind byte   `json:"type_kind"`
}

// FindArrayTypes implements Querier.FindArrayTypes.
func (q *DBQuerier) FindArrayTypes(ctx context.Context, oids []uint32) ([]FindArrayTypesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindArrayTypes")
	rows, err := q.conn.Query(ctx, findArrayTypesSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindArrayTypes: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*uint32)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*string)(nil))
	plan2 := planScan(pgtype.TextCodec{}, fds[2], (*uint32)(nil))
	plan3 := planScan(pgtype.TextCodec{}, fds[3], (*byte)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindArrayTypesRow, error) {
		vals := row.RawValues()
		var item FindArrayTypesRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindArrayTypes.oid: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan FindArrayTypes.type_name: %w", err)
		}
		if err := plan2.Scan(vals[2], &item); err != nil {
			return item, fmt.Errorf("scan FindArrayTypes.elem_oid: %w", err)
		}
		if err := plan3.Scan(vals[3], &item); err != nil {
			return item, fmt.Errorf("scan FindArrayTypes.type_kind: %w", err)
		}
		return item, nil
	})
}

const findCompositeTypesSQL = `WITH table_cols AS (
  SELECT
    cls.relname                                         AS table_name,
    cls.oid                                             AS table_oid,
    array_agg(attr.attname::text ORDER BY attr.attnum)  AS col_names,
    array_agg(attr.atttypid::int8 ORDER BY attr.attnum) AS col_oids,
    array_agg(attr.attnum::int8 ORDER BY attr.attnum)   AS col_orders,
    array_agg(attr.attnotnull ORDER BY attr.attnum)     AS col_not_nulls,
    array_agg(typ.typname::text ORDER BY attr.attnum)   AS col_type_names
  FROM pg_attribute attr
    JOIN pg_class cls ON attr.attrelid = cls.oid
    JOIN pg_type typ ON typ.oid = attr.atttypid
  WHERE attr.attnum > 0 -- Postgres represents system columns with attnum <= 0
    AND NOT attr.attisdropped
  GROUP BY cls.relname, cls.oid
)
SELECT
  typ.typname::text AS table_type_name,
  typ.oid           AS table_type_oid,
  table_name,
  col_names,
  col_oids,
  col_orders,
  col_not_nulls,
  col_type_names
FROM pg_type typ
  JOIN table_cols cols ON typ.typrelid = cols.table_oid
WHERE typ.oid = ANY ($1::oid[])
  AND typ.typtype = 'c';`

type FindCompositeTypesRow struct {
	TableTypeName string                 `json:"table_type_name"`
	TableTypeOID  uint32                 `json:"table_type_oid"`
	TableName     string                 `json:"table_name"`
	ColNames      []string               `json:"col_names"`
	ColOIDs       []int                  `json:"col_oids"`
	ColOrders     []int                  `json:"col_orders"`
	ColNotNulls   pgtype.FlatArray[bool] `json:"col_not_nulls"`
	ColTypeNames  []string               `json:"col_type_names"`
}

// FindCompositeTypes implements Querier.FindCompositeTypes.
func (q *DBQuerier) FindCompositeTypes(ctx context.Context, oids []uint32) ([]FindCompositeTypesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindCompositeTypes")
	rows, err := q.conn.Query(ctx, findCompositeTypesSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindCompositeTypes: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*string)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*uint32)(nil))
	plan2 := planScan(pgtype.TextCodec{}, fds[2], (*string)(nil))
	plan3 := planScan(pgtype.TextCodec{}, fds[3], (*[]string)(nil))
	plan4 := planScan(pgtype.TextCodec{}, fds[4], (*[]int)(nil))
	plan5 := planScan(pgtype.TextCodec{}, fds[5], (*[]int)(nil))
	plan6 := planScan(pgtype.TextCodec{}, fds[6], (*pgtype.FlatArray[bool])(nil))
	plan7 := planScan(pgtype.TextCodec{}, fds[7], (*[]string)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindCompositeTypesRow, error) {
		vals := row.RawValues()
		var item FindCompositeTypesRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.table_type_name: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.table_type_oid: %w", err)
		}
		if err := plan2.Scan(vals[2], &item); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.table_name: %w", err)
		}
		if err := plan3.Scan(vals[3], &item.ColNames); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.col_names: %w", err)
		}
		if err := plan4.Scan(vals[4], &item.ColOIDs); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.col_oids: %w", err)
		}
		if err := plan5.Scan(vals[5], &item.ColOrders); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.col_orders: %w", err)
		}
		if err := plan6.Scan(vals[6], &item); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.col_not_nulls: %w", err)
		}
		if err := plan7.Scan(vals[7], &item.ColTypeNames); err != nil {
			return item, fmt.Errorf("scan FindCompositeTypes.col_type_names: %w", err)
		}
		return item, nil
	})
}

const findDescendantOIDsSQL = `WITH RECURSIVE oid_descs(oid) AS (
  -- Base case.
  SELECT oid
  FROM unnest($1::oid[]) AS t(oid)
  UNION
  -- Recursive case.
  SELECT oid
  FROM (
    WITH all_oids AS (SELECT oid FROM oid_descs)
    -- All composite children.
    SELECT attr.atttypid AS oid
    FROM pg_type typ
      JOIN pg_class cls ON typ.oid = cls.reltype
      JOIN pg_attribute attr ON attr.attrelid = cls.oid
      JOIN all_oids od ON typ.oid = od.oid
    WHERE attr.attnum > 0 -- Postgres represents system columns with attnum <= 0
      AND NOT attr.attisdropped
    UNION
    -- All array elements.
    SELECT elem_typ.oid
    FROM pg_type arr_typ
      JOIN pg_type elem_typ ON arr_typ.typelem = elem_typ.oid
      JOIN all_oids od ON arr_typ.oid = od.oid
  ) t
)
SELECT oid
FROM oid_descs;`

// FindDescendantOIDs implements Querier.FindDescendantOIDs.
func (q *DBQuerier) FindDescendantOIDs(ctx context.Context, oids []uint32) ([]uint32, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindDescendantOIDs")
	rows, err := q.conn.Query(ctx, findDescendantOIDsSQL, oids)
	if err != nil {
		return nil, fmt.Errorf("query FindDescendantOIDs: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*uint32)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (uint32, error) {
		vals := row.RawValues()
		var item uint32
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindDescendantOIDs.oid: %w", err)
		}
		return item, nil
	})
}

const findOIDByNameSQL = `SELECT oid
FROM pg_type
WHERE typname::text = $1
ORDER BY oid DESC
LIMIT 1;`

// FindOIDByName implements Querier.FindOIDByName.
func (q *DBQuerier) FindOIDByName(ctx context.Context, name string) (uint32, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOIDByName")
	rows, err := q.conn.Query(ctx, findOIDByNameSQL, name)
	if err != nil {
		return 0, fmt.Errorf("query FindOIDByName: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*uint32)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (uint32, error) {
		vals := row.RawValues()
		var item uint32
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindOIDByName.oid: %w", err)
		}
		return item, nil
	})
}

const findOIDNameSQL = `SELECT typname AS name
FROM pg_type
WHERE oid = $1;`

// FindOIDName implements Querier.FindOIDName.
func (q *DBQuerier) FindOIDName(ctx context.Context, oid uint32) (string, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOIDName")
	rows, err := q.conn.Query(ctx, findOIDNameSQL, oid)
	if err != nil {
		return "", fmt.Errorf("query FindOIDName: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*string)(nil))

	return pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (string, error) {
		vals := row.RawValues()
		var item string
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindOIDName.name: %w", err)
		}
		return item, nil
	})
}

const findOIDNamesSQL = `SELECT oid, typname AS name, typtype AS kind
FROM pg_type
WHERE oid = ANY ($1::oid[]);`

type FindOIDNamesRow struct {
	OID  uint32 `json:"oid"`
	Name string `json:"name"`
	Kind byte   `json:"kind"`
}

// FindOIDNames implements Querier.FindOIDNames.
func (q *DBQuerier) FindOIDNames(ctx context.Context, oid []uint32) ([]FindOIDNamesRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindOIDNames")
	rows, err := q.conn.Query(ctx, findOIDNamesSQL, oid)
	if err != nil {
		return nil, fmt.Errorf("query FindOIDNames: %w", err)
	}
	fds := rows.FieldDescriptions()
	plan0 := planScan(pgtype.TextCodec{}, fds[0], (*uint32)(nil))
	plan1 := planScan(pgtype.TextCodec{}, fds[1], (*string)(nil))
	plan2 := planScan(pgtype.TextCodec{}, fds[2], (*byte)(nil))

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (FindOIDNamesRow, error) {
		vals := row.RawValues()
		var item FindOIDNamesRow
		if err := plan0.Scan(vals[0], &item); err != nil {
			return item, fmt.Errorf("scan FindOIDNames.oid: %w", err)
		}
		if err := plan1.Scan(vals[1], &item); err != nil {
			return item, fmt.Errorf("scan FindOIDNames.name: %w", err)
		}
		if err := plan2.Scan(vals[2], &item); err != nil {
			return item, fmt.Errorf("scan FindOIDNames.kind: %w", err)
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
