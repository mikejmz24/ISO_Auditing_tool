package testutils

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*sql.Rows), argsMock.Error(1)
}

func (m *MockDatabase) Exec(query string, args ...any) (sql.Result, error) {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(sql.Result), argsMock.Error(1)
}
