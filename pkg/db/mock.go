//nolint:all // файл содержит моки для тестов, проверки линтеров не требуются
package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockPool struct {
	mock.Mock
}

func (m *MockPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPool) Close() {
	m.Called()
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgx.Row)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgx.Rows), called.Error(1)
}

func (m *MockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	called := m.Called(ctx)
	return called.Get(0).(pgx.Tx), called.Error(1)
}

func (m *MockPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	called := m.Called(ctx, sql, arguments)
	return called.Get(0).(pgconn.CommandTag), called.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRow) Close() {
	m.Called()
}

func (m *MockRow) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRow) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func (m *MockRow) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRow) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

func (m *MockRow) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRow) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func (m *MockTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgconn.CommandTag), called.Error(1)
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockTx) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgx.Rows), called.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	called := m.Called(ctx, sql, args)
	return called.Get(0).(pgx.Row)
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier,
	columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	return args.Get(0).(pgx.BatchResults)
}

type MockBatchResults struct {
	mock.Mock
}

func (m *MockBatchResults) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *MockBatchResults) Query() (pgx.Rows, error) {
	args := m.Called()
	return args.Get(0).(pgx.Rows), args.Error(1)
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	args := m.Called()
	return args.Get(0).(pgx.Row)
}

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Close() {
	m.Called()
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}
