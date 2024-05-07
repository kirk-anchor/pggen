// Code generated by pggen. DO NOT EDIT.

package device

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type QueryName struct{}

// Querier is a typesafe Go interface backed by SQL queries.
type Querier interface {
	FindDevicesByUser(ctx context.Context, id int) ([]FindDevicesByUserRow, error)

	CompositeUser(ctx context.Context) ([]CompositeUserRow, error)

	CompositeUserOne(ctx context.Context) (User, error)

	CompositeUserOneTwoCols(ctx context.Context) (CompositeUserOneTwoColsRow, error)

	CompositeUserMany(ctx context.Context) ([]User, error)

	InsertUser(ctx context.Context, userID int, name string) (pgconn.CommandTag, error)

	InsertDevice(ctx context.Context, mac net.HardwareAddr, owner int) (pgconn.CommandTag, error)
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

// User represents the Postgres composite type "user".
type User struct {
	ID   *int    `json:"id"`
	Name *string `json:"name"`
}

// DeviceType represents the Postgres enum "device_type".
type DeviceType string

const (
	DeviceTypeUndefined DeviceType = "undefined"
	DeviceTypePhone     DeviceType = "phone"
	DeviceTypeLaptop    DeviceType = "laptop"
	DeviceTypeIpad      DeviceType = "ipad"
	DeviceTypeDesktop   DeviceType = "desktop"
	DeviceTypeIot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }

const findDevicesByUserSQL = `SELECT
  id,
  name,
  (SELECT array_agg(mac) FROM device WHERE owner = id) AS mac_addrs
FROM "user"
WHERE id = $1;`

type FindDevicesByUserRow struct {
	ID       int                                `json:"id"`
	Name     string                             `json:"name"`
	MacAddrs pgtype.FlatArray[net.HardwareAddr] `json:"mac_addrs"`
}

// FindDevicesByUser implements Querier.FindDevicesByUser.
func (q *DBQuerier) FindDevicesByUser(ctx context.Context, id int) ([]FindDevicesByUserRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "FindDevicesByUser")
	rows, err := q.conn.Query(ctx, findDevicesByUserSQL, id)
	if err != nil {
		return nil, fmt.Errorf("query FindDevicesByUser: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[FindDevicesByUserRow])
}

const compositeUserSQL = `SELECT
  d.mac,
  d.type,
  ROW (u.id, u.name)::"user" AS "user"
FROM device d
  LEFT JOIN "user" u ON u.id = d.owner;`

type CompositeUserRow struct {
	Mac  net.HardwareAddr `json:"mac"`
	Type DeviceType       `json:"type"`
	User User             `json:"user"`
}

// CompositeUser implements Querier.CompositeUser.
func (q *DBQuerier) CompositeUser(ctx context.Context) ([]CompositeUserRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CompositeUser")
	rows, err := q.conn.Query(ctx, compositeUserSQL)
	if err != nil {
		return nil, fmt.Errorf("query CompositeUser: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[CompositeUserRow])
}

const compositeUserOneSQL = `SELECT ROW (15, 'qux')::"user" AS "user";`

// CompositeUserOne implements Querier.CompositeUserOne.
func (q *DBQuerier) CompositeUserOne(ctx context.Context) (User, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CompositeUserOne")
	rows, err := q.conn.Query(ctx, compositeUserOneSQL)
	if err != nil {
		return User{}, fmt.Errorf("query CompositeUserOne: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowTo[User])
}

const compositeUserOneTwoColsSQL = `SELECT 1 AS num, ROW (15, 'qux')::"user" AS "user";`

type CompositeUserOneTwoColsRow struct {
	Num  int32 `json:"num"`
	User User  `json:"user"`
}

// CompositeUserOneTwoCols implements Querier.CompositeUserOneTwoCols.
func (q *DBQuerier) CompositeUserOneTwoCols(ctx context.Context) (CompositeUserOneTwoColsRow, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CompositeUserOneTwoCols")
	rows, err := q.conn.Query(ctx, compositeUserOneTwoColsSQL)
	if err != nil {
		return CompositeUserOneTwoColsRow{}, fmt.Errorf("query CompositeUserOneTwoCols: %w", err)
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CompositeUserOneTwoColsRow])
}

const compositeUserManySQL = `SELECT ROW (15, 'qux')::"user" AS "user";`

// CompositeUserMany implements Querier.CompositeUserMany.
func (q *DBQuerier) CompositeUserMany(ctx context.Context) ([]User, error) {
	ctx = context.WithValue(ctx, QueryName{}, "CompositeUserMany")
	rows, err := q.conn.Query(ctx, compositeUserManySQL)
	if err != nil {
		return nil, fmt.Errorf("query CompositeUserMany: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowTo[User])
}

const insertUserSQL = `INSERT INTO "user" (id, name)
VALUES ($1, $2);`

// InsertUser implements Querier.InsertUser.
func (q *DBQuerier) InsertUser(ctx context.Context, userID int, name string) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertUser")
	cmdTag, err := q.conn.Exec(ctx, insertUserSQL, userID, name)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertUser: %w", err)
	}
	return cmdTag, err
}

const insertDeviceSQL = `INSERT INTO device (mac, owner)
VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac net.HardwareAddr, owner int) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, QueryName{}, "InsertDevice")
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, owner)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
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
