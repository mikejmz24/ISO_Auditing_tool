// tests/unit/repositories/iso_standard_repository_test.go
package repositories_test

import (
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func TestGetAllISOStandards(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	expectedStandards := []types.ISOStandard{
		{ID: 1, Name: "ISO 9001"},
	}

	mockRepo.On("GetAllISOStandards").Return(expectedStandards, nil)

	standards, err := mockRepo.GetAllISOStandards()
	assert.NoError(t, err)
	assert.Equal(t, expectedStandards, standards)

	mockRepo.AssertExpectations(t)
}

func TestCreateISOStandard(t *testing.T) {
	mockRepo := new(testutils.MockIsoStandardRepository)
	newStandard := types.ISOStandard{Name: "ISO 9001"}
	expectedID := int64(1)

	mockRepo.On("CreateISOStandard", newStandard).Return(expectedID, nil)

	id, err := mockRepo.CreateISOStandard(newStandard)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	mockRepo.AssertExpectations(t)
}
