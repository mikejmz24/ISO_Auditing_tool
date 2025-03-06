package testutils

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDraftRepository struct {
	mock.Mock
}

// func (m *MockIsoStandardRepository) GetAllISOStandards() ([]types.ISOStandard, error) {
// 	args := m.Called()
// 	// Assert the mocked response is of the expected type
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]types.ISOStandard), args.Error(1)
// }
//
// func (m *MockIsoStandardRepository) GetISOStandardByID(id int64) (types.ISOStandard, error) {
// 	args := m.Called(id)
// 	return args.Get(0).(types.ISOStandard), args.Error(1)
// }

func (m *MockDraftRepository) Create(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

//	func (m *MockIsoStandardRepository) UpdateISOStandard(isoStandard types.ISOStandard) error {
//		args := m.Called(isoStandard)
//		return args.Error(0)
//	}
//
//	func (m *MockIsoStandardRepository) DeleteISOStandard(id int64) error {
//		args := m.Called(id)
//		return args.Error(0)
//	}
func (m *MockDraftRepository) Reset() {
	m.ExpectedCalls = nil
	m.Calls = nil
}
