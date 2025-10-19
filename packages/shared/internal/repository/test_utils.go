package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockDBTX is a mock implementation of DBTX interface
type MockDBTX struct {
	mock.Mock
}

func (m *MockDBTX) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsMock := m.Called(ctx, query, args)
	if argsMock.Get(0) == nil {
		return pgconn.CommandTag{}, argsMock.Error(1)
	}
	return argsMock.Get(0).(pgconn.CommandTag), argsMock.Error(1)
}

func (m *MockDBTX) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	argsMock := m.Called(ctx, query, args)
	if argsMock.Get(0) == nil {
		return nil, argsMock.Error(1)
	}
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

func (m *MockDBTX) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	argsMock := m.Called(ctx, query, args)
	return argsMock.Get(0).(pgx.Row)
}

// MockTx is a mock implementation of pgx.Tx interface
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	argsMock := m.Called(ctx, sql, arguments)
	if argsMock.Get(0) == nil {
		return pgconn.CommandTag{}, argsMock.Error(1)
	}
	return argsMock.Get(0).(pgconn.CommandTag), argsMock.Error(1)
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	argsMock := m.Called(ctx, sql, args)
	if argsMock.Get(0) == nil {
		return nil, argsMock.Error(1)
	}
	return argsMock.Get(0).(pgx.Rows), argsMock.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.Row)
}

// Implement other required pgx.Tx methods as no-ops for testing
func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (m *MockTx) Commit(ctx context.Context) error {
	return nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	return nil
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *MockTx) Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *MockTx) Conn() *pgx.Conn {
	return nil
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
