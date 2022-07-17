// Code generated by pggen. DO NOT EDIT.

package slices

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"time"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	GetBools(ctx context.Context, data []bool) ([]bool, error)
	// GetBoolsBatch enqueues a GetBools query into batch to be executed
	// later by the batch.
	GetBoolsBatch(batch genericBatch, data []bool)
	// GetBoolsScan scans the result of an executed GetBoolsBatch query.
	GetBoolsScan(results pgx.BatchResults) ([]bool, error)

	GetOneTimestamp(ctx context.Context, data *time.Time) (*time.Time, error)
	// GetOneTimestampBatch enqueues a GetOneTimestamp query into batch to be executed
	// later by the batch.
	GetOneTimestampBatch(batch genericBatch, data *time.Time)
	// GetOneTimestampScan scans the result of an executed GetOneTimestampBatch query.
	GetOneTimestampScan(results pgx.BatchResults) (*time.Time, error)

	GetManyTimestamptzs(ctx context.Context, data []time.Time) ([]*time.Time, error)
	// GetManyTimestamptzsBatch enqueues a GetManyTimestamptzs query into batch to be executed
	// later by the batch.
	GetManyTimestamptzsBatch(batch genericBatch, data []time.Time)
	// GetManyTimestamptzsScan scans the result of an executed GetManyTimestamptzsBatch query.
	GetManyTimestamptzsScan(results pgx.BatchResults) ([]*time.Time, error)

	GetManyTimestamps(ctx context.Context, data []*time.Time) ([]*time.Time, error)
	// GetManyTimestampsBatch enqueues a GetManyTimestamps query into batch to be executed
	// later by the batch.
	GetManyTimestampsBatch(batch genericBatch, data []*time.Time)
	// GetManyTimestampsScan scans the result of an executed GetManyTimestampsBatch query.
	GetManyTimestampsScan(results pgx.BatchResults) ([]*time.Time, error)
}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// genericBatch batches queries to send in a single network request to a
// Postgres server. This is usually backed by *pgx.Batch.
type genericBatch interface {
	// Queue queues a query to batch b. query can be an SQL query or the name of a
	// prepared statement. See Queue on *pgx.Batch.
	Queue(query string, arguments ...interface{})
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return NewQuerierConfig(conn, QuerierConfig{})
}

type QuerierConfig struct {
	// DataTypes contains pgtype.Value to use for encoding and decoding instead
	// of pggen-generated pgtype.ValueTranscoder.
	//
	// If OIDs are available for an input parameter type and all of its
	// transitive dependencies, pggen will use the binary encoding format for
	// the input parameter.
	DataTypes []pgtype.DataType
}

// NewQuerierConfig creates a DBQuerier that implements Querier with the given
// config. conn is typically *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerierConfig(conn genericConn, cfg QuerierConfig) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver(cfg.DataTypes)}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, getBoolsSQL, getBoolsSQL); err != nil {
		return fmt.Errorf("prepare query 'GetBools': %w", err)
	}
	if _, err := p.Prepare(ctx, getOneTimestampSQL, getOneTimestampSQL); err != nil {
		return fmt.Errorf("prepare query 'GetOneTimestamp': %w", err)
	}
	if _, err := p.Prepare(ctx, getManyTimestamptzsSQL, getManyTimestamptzsSQL); err != nil {
		return fmt.Errorf("prepare query 'GetManyTimestamptzs': %w", err)
	}
	if _, err := p.Prepare(ctx, getManyTimestampsSQL, getManyTimestampsSQL); err != nil {
		return fmt.Errorf("prepare query 'GetManyTimestamps': %w", err)
	}
	return nil
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver(types []pgtype.DataType) *typeResolver {
	ci := pgtype.NewConnInfo()
	for _, typ := range types {
		if txt, ok := typ.Value.(textPreferrer); ok && typ.OID != unknownOID {
			typ.Value = txt.ValueTranscoder
		}
		ci.RegisterDataType(typ)
	}
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

type compositeField struct {
	name       string                 // name of the field
	typeName   string                 // Postgres type name
	defaultVal pgtype.ValueTranscoder // default value to use
}

func (tr *typeResolver) newCompositeValue(name string, fields ...compositeField) pgtype.ValueTranscoder {
	if _, val, ok := tr.findValue(name); ok {
		return val
	}
	fs := make([]pgtype.CompositeTypeField, len(fields))
	vals := make([]pgtype.ValueTranscoder, len(fields))
	isBinaryOk := true
	for i, field := range fields {
		oid, val, ok := tr.findValue(field.typeName)
		if !ok {
			oid = unknownOID
			val = field.defaultVal
		}
		isBinaryOk = isBinaryOk && oid != unknownOID
		fs[i] = pgtype.CompositeTypeField{Name: field.name, OID: oid}
		vals[i] = val
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	typ, _ := pgtype.NewCompositeTypeValues(name, fs, vals)
	if !isBinaryOk {
		return textPreferrer{typ, name}
	}
	return typ
}

func (tr *typeResolver) newArrayValue(name, elemName string, defaultVal func() pgtype.ValueTranscoder) pgtype.ValueTranscoder {
	if _, val, ok := tr.findValue(name); ok {
		return val
	}
	elemOID, elemVal, ok := tr.findValue(elemName)
	elemValFunc := func() pgtype.ValueTranscoder {
		return pgtype.NewValue(elemVal).(pgtype.ValueTranscoder)
	}
	if !ok {
		elemOID = unknownOID
		elemValFunc = defaultVal
	}
	typ := pgtype.NewArrayType(name, elemOID, elemValFunc)
	if elemOID == unknownOID {
		return textPreferrer{typ, name}
	}
	return typ
}

// newboolArrayRaw returns all elements for the Postgres array type '_bool'
// as a slice of interface{} for use with the pgtype.Value Set method.
func (tr *typeResolver) newboolArrayRaw(vs []bool) []interface{} {
	elems := make([]interface{}, len(vs))
	for i, v := range vs {
		elems[i] = v
	}
	return elems
}

const getBoolsSQL = `SELECT $1::boolean[];`

// GetBools implements Querier.GetBools.
func (q *DBQuerier) GetBools(ctx context.Context, data []bool) ([]bool, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetBools")
	row := q.conn.QueryRow(ctx, getBoolsSQL, data)
	item := []bool{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetBools: %w", err)
	}
	return item, nil
}

// GetBoolsBatch implements Querier.GetBoolsBatch.
func (q *DBQuerier) GetBoolsBatch(batch genericBatch, data []bool) {
	batch.Queue(getBoolsSQL, data)
}

// GetBoolsScan implements Querier.GetBoolsScan.
func (q *DBQuerier) GetBoolsScan(results pgx.BatchResults) ([]bool, error) {
	row := results.QueryRow()
	item := []bool{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetBoolsBatch row: %w", err)
	}
	return item, nil
}

const getOneTimestampSQL = `SELECT $1::timestamp;`

// GetOneTimestamp implements Querier.GetOneTimestamp.
func (q *DBQuerier) GetOneTimestamp(ctx context.Context, data *time.Time) (*time.Time, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetOneTimestamp")
	row := q.conn.QueryRow(ctx, getOneTimestampSQL, data)
	var item time.Time
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("query GetOneTimestamp: %w", err)
	}
	return &item, nil
}

// GetOneTimestampBatch implements Querier.GetOneTimestampBatch.
func (q *DBQuerier) GetOneTimestampBatch(batch genericBatch, data *time.Time) {
	batch.Queue(getOneTimestampSQL, data)
}

// GetOneTimestampScan implements Querier.GetOneTimestampScan.
func (q *DBQuerier) GetOneTimestampScan(results pgx.BatchResults) (*time.Time, error) {
	row := results.QueryRow()
	var item time.Time
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("scan GetOneTimestampBatch row: %w", err)
	}
	return &item, nil
}

const getManyTimestamptzsSQL = `SELECT *
FROM unnest($1::timestamptz[]);`

// GetManyTimestamptzs implements Querier.GetManyTimestamptzs.
func (q *DBQuerier) GetManyTimestamptzs(ctx context.Context, data []time.Time) ([]*time.Time, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetManyTimestamptzs")
	rows, err := q.conn.Query(ctx, getManyTimestamptzsSQL, data)
	if err != nil {
		return nil, fmt.Errorf("query GetManyTimestamptzs: %w", err)
	}
	defer rows.Close()
	items := []*time.Time{}
	for rows.Next() {
		var item time.Time
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan GetManyTimestamptzs row: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetManyTimestamptzs rows: %w", err)
	}
	return items, err
}

// GetManyTimestamptzsBatch implements Querier.GetManyTimestamptzsBatch.
func (q *DBQuerier) GetManyTimestamptzsBatch(batch genericBatch, data []time.Time) {
	batch.Queue(getManyTimestamptzsSQL, data)
}

// GetManyTimestamptzsScan implements Querier.GetManyTimestamptzsScan.
func (q *DBQuerier) GetManyTimestamptzsScan(results pgx.BatchResults) ([]*time.Time, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query GetManyTimestamptzsBatch: %w", err)
	}
	defer rows.Close()
	items := []*time.Time{}
	for rows.Next() {
		var item time.Time
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan GetManyTimestamptzsBatch row: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetManyTimestamptzsBatch rows: %w", err)
	}
	return items, err
}

const getManyTimestampsSQL = `SELECT *
FROM unnest($1::timestamp[]);`

// GetManyTimestamps implements Querier.GetManyTimestamps.
func (q *DBQuerier) GetManyTimestamps(ctx context.Context, data []*time.Time) ([]*time.Time, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetManyTimestamps")
	rows, err := q.conn.Query(ctx, getManyTimestampsSQL, data)
	if err != nil {
		return nil, fmt.Errorf("query GetManyTimestamps: %w", err)
	}
	defer rows.Close()
	items := []*time.Time{}
	for rows.Next() {
		var item time.Time
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan GetManyTimestamps row: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetManyTimestamps rows: %w", err)
	}
	return items, err
}

// GetManyTimestampsBatch implements Querier.GetManyTimestampsBatch.
func (q *DBQuerier) GetManyTimestampsBatch(batch genericBatch, data []*time.Time) {
	batch.Queue(getManyTimestampsSQL, data)
}

// GetManyTimestampsScan implements Querier.GetManyTimestampsScan.
func (q *DBQuerier) GetManyTimestampsScan(results pgx.BatchResults) ([]*time.Time, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query GetManyTimestampsBatch: %w", err)
	}
	defer rows.Close()
	items := []*time.Time{}
	for rows.Next() {
		var item time.Time
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan GetManyTimestampsBatch row: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetManyTimestampsBatch rows: %w", err)
	}
	return items, err
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs if the results from QueryAllDataTypes aren't passed to
// NewQuerierConfig.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0