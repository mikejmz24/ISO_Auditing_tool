package testutils

import (
	"ISO_Auditing_Tool/pkg/types"
	"github.com/stretchr/testify/mock"
)

type MockIsoStandardRepository struct {
	mock.Mock
}

func (m *MockIsoStandardRepository) GetAllISOStandards() ([]types.ISOStandard, error) {
	args := m.Called()
	return args.Get(0).([]types.ISOStandard), args.Error(1)
}

func (m *MockIsoStandardRepository) GetISOStandardByID(id int) (types.ISOStandard, error) {
	args := m.Called(id)
	return args.Get(0).(types.ISOStandard), args.Error(1)
}

func (m *MockIsoStandardRepository) CreateISOStandard(standard types.ISOStandard) (int64, error) {
	args := m.Called(standard)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockIsoStandardRepository) UpdateISOStandard(standard types.ISOStandard) error {
	args := m.Called(standard)
	return args.Error(0)
}

func (m *MockIsoStandardRepository) DeleteISOStandard(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
