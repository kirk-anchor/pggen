// Code generated by pggen. DO NOT EDIT.

package pgcrypto

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	CreateUser(ctx context.Context, email string, password string) (pgconn.CommandTag, error)
	// CreateUserBatch enqueues a CreateUser query into batch to be executed
	// later by the batch.
	CreateUserBatch(batch genericBatch, email string, password string)
	// CreateUserScan scans the result of an executed CreateUserBatch query.
	CreateUserScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	FindUser(ctx context.Context, email string) (FindUserRow, error)
	// FindUserBatch enqueues a FindUser query into batch to be executed
	// later by the batch.
	FindUserBatch(batch genericBatch, email string)
	// FindUserScan scans the result of an executed FindUserBatch query.
	FindUserScan(results pgx.BatchResults) (FindUserRow, error)
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
	if _, err := p.Prepare(ctx, createUserSQL, createUserSQL); err != nil {
		return fmt.Errorf("prepare query 'CreateUser': %w", err)
	}
	if _, err := p.Prepare(ctx, findUserSQL, findUserSQL); err != nil {
		return fmt.Errorf("prepare query 'FindUser': %w", err)
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

const createUserSQL = `INSERT INTO "user" (email, pass)
VALUES ($1, crypt($2, gen_salt('bf')));`

// CreateUser implements Querier.CreateUser.
func (q *DBQuerier) CreateUser(ctx context.Context, email string, password string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CreateUser")
	cmdTag, err := q.conn.Exec(ctx, createUserSQL, email, password)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query CreateUser: %w", err)
	}
	return cmdTag, err
}

// CreateUserBatch implements Querier.CreateUserBatch.
func (q *DBQuerier) CreateUserBatch(batch genericBatch, email string, password string) {
	batch.Queue(createUserSQL, email, password)
}

// CreateUserScan implements Querier.CreateUserScan.
func (q *DBQuerier) CreateUserScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec CreateUserBatch: %w", err)
	}
	return cmdTag, err
}

const findUserSQL = `SELECT email, pass from "user"
where email = $1;`

type FindUserRow struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

// FindUser implements Querier.FindUser.
func (q *DBQuerier) FindUser(ctx context.Context, email string) (FindUserRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindUser")
	row := q.conn.QueryRow(ctx, findUserSQL, email)
	var item FindUserRow
	if err := row.Scan(&item.Email, &item.Pass); err != nil {
		return item, fmt.Errorf("query FindUser: %w", err)
	}
	return item, nil
}

// FindUserBatch implements Querier.FindUserBatch.
func (q *DBQuerier) FindUserBatch(batch genericBatch, email string) {
	batch.Queue(findUserSQL, email)
}

// FindUserScan implements Querier.FindUserScan.
func (q *DBQuerier) FindUserScan(results pgx.BatchResults) (FindUserRow, error) {
	row := results.QueryRow()
	var item FindUserRow
	if err := row.Scan(&item.Email, &item.Pass); err != nil {
		return item, fmt.Errorf("scan FindUserBatch row: %w", err)
	}
	return item, nil
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
	return textPreferrer{ValueTranscoder: pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), typeName: t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
