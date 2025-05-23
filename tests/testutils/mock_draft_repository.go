package testutils

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDraftRepository struct {
	mock.Mock
}

func (m *MockDraftRepository) CreateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) GetDraftByID(ctx context.Context, id int) (types.Draft, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) GetByID(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) Delete(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) GetDraftsByTypeAndObject(ctx context.Context, typeID, objectID int) ([]types.Draft, error) {
	args := m.Called(ctx, typeID, objectID)
	return args.Get(0).([]types.Draft), args.Error(1)
}

func (m *MockDraftRepository) Create(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) Reset() {
	m.ExpectedCalls = nil
	m.Calls = nil
}
