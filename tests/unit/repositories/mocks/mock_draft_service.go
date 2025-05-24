package mocks

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockDraftService struct {
	mock.Mock
}

func (m *MockDraftService) Create(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftService) GetByID(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftService) Update(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftService) Reset() {
	m.ExpectedCalls = nil
	m.Calls = nil
}
