package mocks

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

func (m *MockDraftRepository) GetDraftByID(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) UpdateDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) DeleteDraft(ctx context.Context, draft types.Draft) (types.Draft, error) {
	args := m.Called(ctx, draft)
	return args.Get(0).(types.Draft), args.Error(1)
}

func (m *MockDraftRepository) GetDraftsByTypeAndObject(ctx context.Context, typeID, objectID int) ([]types.Draft, error) {
	args := m.Called(ctx, typeID, objectID)
	return args.Get(0).([]types.Draft), args.Error(1)
}

func (m *MockDraftRepository) UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error) {
	args := m.Called(ctx, requirement)
	return args.Get(0).(types.Requirement), args.Error(1)
}

// Reset clears all expectations and calls
func (m *MockDraftRepository) Reset() {
	m.ExpectedCalls = nil
	m.Calls = nil
}
