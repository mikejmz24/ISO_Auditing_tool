package controller_test

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewDraftController(t *testing.T) {
	// Arrange
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Act
	repo, err := repositories.NewDraftRepository(db)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestDraftController_Success(t *testing.T) {
	// Common setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := repositories.NewDraftRepository(db)
	assert.NoError(t, err)

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	data := json.RawMessage(`{"field1": "value1", "field2": "value2"}`)
	diff := json.RawMessage(`{"changed": "content"}`)

	tests := []struct {
		name       string
		draft      types.Draft
		mockSetup  func()
		expectedID int
	}{
		{
			name: "Create",
			draft: types.Draft{
				TypeID:          1,
				ObjectID:        2,
				StatusID:        3,
				Version:         1,
				Data:            data,
				Diff:            diff,
				UserID:          42,
				ApproverID:      0,
				ApprovalComment: "",
				PublishError:    "",
				CreatedAt:       now,
				UpdatedAt:       now,
				ExpiresAt:       now.Add(24 * time.Hour),
			},
			mockSetup: func() {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up test case
			tc.mockSetup()

			// Act
			result, err := repo.Create(ctx, tc.draft)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedID, result.ID)

			// Verify expectations for this test case
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestDraftController_Errors(t *testing.T) {
	// Common setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo, err := repositories.NewDraftRepository(db)
	assert.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name          string
		draft         types.Draft
		mockSetup     func()
		expectedError string
	}{
		{
			name: "An execution error on Create returns failed to create draft",
			draft: types.Draft{
				TypeID:   1,
				ObjectID: 2,
				UserID:   42,
			},
			mockSetup: func() {
				execError := errors.New("execution error")
				mock.ExpectExec("").WillReturnError(execError)
			},
			expectedError: "failed to create draft",
		},
		{
			name: "SQL last insert error on Create returns failed to get last insert ID",
			draft: types.Draft{
				TypeID:   1,
				ObjectID: 2,
				UserID:   42,
			},
			mockSetup: func() {
				lastIDError := errors.New("last insert ID error")
				mockResult := sqlmock.NewErrorResult(lastIDError)
				mock.ExpectExec("").WillReturnResult(mockResult)
			},
			expectedError: "failed to get last insert ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up test case
			tc.mockSetup()

			// Act
			result, err := repo.Create(ctx, tc.draft)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, types.Draft{}, result)
			assert.Contains(t, err.Error(), tc.expectedError)

			// Verify expectations for this test case
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
