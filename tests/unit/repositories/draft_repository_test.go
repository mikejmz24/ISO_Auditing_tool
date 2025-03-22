package repositories

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

// DraftRepositorySuccessTestSuite defines the test suite for successful operations
type DraftRepositorySuccessTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo repositories.DraftRepository
}

// DraftRepositoryFailureTestSuite defines the test suite for failure operations
type DraftRepositoryFailureTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo repositories.DraftRepository
}

// SetupTest initializes test dependencies before each test
func (s *DraftRepositorySuccessTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo, err := repositories.NewDraftRepository(db)
	s.NoError(err)

	s.db = db
	s.mock = mock
	s.repo = repo
}

// TearDownTest cleans up resources after each test
func (s *DraftRepositorySuccessTestSuite) TearDownTest() {
	s.db.Close()
}

// SetupTest initializes test dependencies before each test
func (s *DraftRepositoryFailureTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo, err := repositories.NewDraftRepository(db)
	s.NoError(err)

	s.db = db
	s.mock = mock
	s.repo = repo
}

// TearDownTest cleans up resources after each test
func (s *DraftRepositoryFailureTestSuite) TearDownTest() {
	s.db.Close()
}

// TestCreateSuccess tests successful draft creation scenarios
func (s *DraftRepositorySuccessTestSuite) TestCreateSuccess() {
	testCases := []struct {
		name     string
		draft    types.Draft
		expected types.Draft
	}{
		{
			name: "Basic_draft",
			draft: types.Draft{
				TypeID:          1,
				ObjectID:        2,
				StatusID:        3,
				Version:         1,
				Data:            json.RawMessage(`{"field1": "value1", "field2": "value2"}`),
				Diff:            json.RawMessage(`{"changed": "content"}`),
				UserID:          42,
				ApproverID:      0,
				ApprovalComment: "",
				PublishError:    "",
				CreatedAt:       time.Now().UTC().Truncate(time.Second),
				UpdatedAt:       time.Now().UTC().Truncate(time.Second),
				ExpiresAt:       time.Now().UTC().Truncate(time.Second).Add(24 * time.Hour),
			},
			expected: types.Draft{
				ID:              1,
				TypeID:          1,
				ObjectID:        2,
				StatusID:        3,
				Version:         1,
				Data:            json.RawMessage(`{"field1": "value1", "field2": "value2"}`),
				Diff:            json.RawMessage(`{"changed": "content"}`),
				UserID:          42,
				ApproverID:      0,
				ApprovalComment: "",
				PublishError:    "",
			},
		},
		{
			name: "Draft_with_approver_id",
			draft: types.Draft{
				TypeID:          1,
				ObjectID:        2,
				StatusID:        3,
				Version:         2,
				Data:            json.RawMessage(`{"field1": "updated"}`),
				Diff:            json.RawMessage(`{"field1": {"old": "value1", "new": "updated"}}`),
				UserID:          42,
				ApproverID:      10,
				ApprovalComment: "Looks good",
				PublishError:    "",
				CreatedAt:       time.Now().UTC().Truncate(time.Second),
				UpdatedAt:       time.Now().UTC().Truncate(time.Second),
				ExpiresAt:       time.Now().UTC().Truncate(time.Second).Add(24 * time.Hour),
			},
			expected: types.Draft{
				ID:              1,
				TypeID:          1,
				ObjectID:        2,
				StatusID:        3,
				Version:         2,
				Data:            json.RawMessage(`{"field1": "updated"}`),
				Diff:            json.RawMessage(`{"field1": {"old": "value1", "new": "updated"}}`),
				UserID:          42,
				ApproverID:      10,
				ApprovalComment: "Looks good",
				PublishError:    "",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Arrange - set up the mock
			s.mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			// Act
			result, err := s.repo.Create(context.Background(), tc.draft)

			// Assert
			s.NoError(err)
			s.Equal(tc.expected.ID, result.ID)
			s.Equal(tc.expected.TypeID, result.TypeID)
			s.Equal(tc.expected.ObjectID, result.ObjectID)
			s.Equal(tc.expected.StatusID, result.StatusID)
			s.Equal(tc.expected.Version, result.Version)
			s.Equal(tc.expected.UserID, result.UserID)
			s.Equal(string(tc.expected.Data), string(result.Data))
			s.Equal(string(tc.expected.Diff), string(result.Diff))
			s.Equal(tc.expected.ApproverID, result.ApproverID)
			s.Equal(tc.expected.ApprovalComment, result.ApprovalComment)
			s.Equal(tc.expected.PublishError, result.PublishError)

			// Verify expectations were met
			s.NoError(s.mock.ExpectationsWereMet())
		})
	}
}

// TestCreateFailure tests scenarios where draft creation fails
// TODO: Return custom errors
func (s *DraftRepositoryFailureTestSuite) TestCreateFailure() {
	now := time.Now().UTC()
	baseDraft := types.Draft{
		TypeID:          1,
		ObjectID:        2,
		StatusID:        3,
		Version:         1,
		Data:            json.RawMessage(`{"field1": "value1", "field2": "value2"}`),
		Diff:            json.RawMessage(`{"changed": "content"}`),
		UserID:          42,
		ApproverID:      0,
		ApprovalComment: "",
		PublishError:    "",
		CreatedAt:       now,
		UpdatedAt:       now,
		ExpiresAt:       now.Add(24 * time.Hour),
	}

	testCases := []struct {
		name          string
		draft         types.Draft
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
	}{
		{
			name:  "SQL_execution_error_returns_failed_to_create_draft",
			draft: baseDraft,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("").WillReturnError(errors.New("execution error"))
			},
			expectedError: "failed to create draft",
		},
		{
			name:  "SQL_last_insert_error_returns_failed_to_get_last_insert_if",
			draft: baseDraft,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mockResult := sqlmock.NewErrorResult(errors.New("last insert ID error"))
				mock.ExpectExec("").WillReturnResult(mockResult)
			},
			expectedError: "failed to get last insert ID",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Arrange
			tc.mockSetup(s.mock)

			// Act
			result, err := s.repo.Create(context.Background(), tc.draft)

			// Assert
			s.Error(err)
			s.Contains(err.Error(), tc.expectedError)
			s.Equal(types.Draft{}, result)

			// Verify expectations were met
			s.NoError(s.mock.ExpectationsWereMet())
		})
	}
}

func (s *DraftRepositorySuccessTestSuite) TestUpdateSuccess() {
	draftUpdate := types.Draft{
		TypeID:          1,
		ObjectID:        2,
		StatusID:        3,
		Version:         1,
		Data:            json.RawMessage(`{"field1": "value1", "field2": "value2"}`),
		Diff:            json.RawMessage(`{"changed": "content"}`),
		UserID:          42,
		ApproverID:      0,
		ApprovalComment: "",
		PublishError:    "",
		CreatedAt:       time.Now().UTC().Truncate(time.Second),
		UpdatedAt:       time.Now().UTC().Truncate(time.Second),
		ExpiresAt:       time.Now().UTC().Truncate(time.Second).Add(24 * time.Hour),
	}
	testCases := []struct {
		name     string
		draft    types.Draft
		expected types.Draft
	}{
		{
			name:     "Basic_draft",
			draft:    draftUpdate,
			expected: draftUpdate,
		},
		// 	{
		// 		name: "Draft_with_approver_id",
		// 		draft: types.Draft{
		// 			TypeID:          1,
		// 			ObjectID:        2,
		// 			StatusID:        3,
		// 			Version:         2,
		// 			Data:            json.RawMessage(`{"field1": "updated"}`),
		// 			Diff:            json.RawMessage(`{"field1": {"old": "value1", "new": "updated"}}`),
		// 			UserID:          42,
		// 			ApproverID:      10,
		// 			ApprovalComment: "Looks good",
		// 			PublishError:    "",
		// 			CreatedAt:       time.Now().UTC().Truncate(time.Second),
		// 			UpdatedAt:       time.Now().UTC().Truncate(time.Second),
		// 			ExpiresAt:       time.Now().UTC().Truncate(time.Second).Add(24 * time.Hour),
		// 		},
		// 		expected: types.Draft{
		// 			ID:              1,
		// 			TypeID:          1,
		// 			ObjectID:        2,
		// 			StatusID:        3,
		// 			Version:         2,
		// 			Data:            json.RawMessage(`{"field1": "updated"}`),
		// 			Diff:            json.RawMessage(`{"field1": {"old": "value1", "new": "updated"}}`),
		// 			UserID:          42,
		// 			ApproverID:      10,
		// 			ApprovalComment: "Looks good",
		// 			PublishError:    "",
		// 		},
		// 	},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Arrange - set up the mock
			s.mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))

			// Act
			result, err := s.repo.Update(context.Background(), tc.draft)

			// Assert
			s.NoError(err)
			s.Equal(tc.expected.ID, result.ID)
			s.Equal(tc.expected.TypeID, result.TypeID)
			s.Equal(tc.expected.ObjectID, result.ObjectID)
			s.Equal(tc.expected.StatusID, result.StatusID)
			s.Equal(tc.expected.Version, result.Version)
			s.Equal(tc.expected.UserID, result.UserID)
			s.Equal(string(tc.expected.Data), string(result.Data))
			s.Equal(string(tc.expected.Diff), string(result.Diff))
			s.Equal(tc.expected.ApproverID, result.ApproverID)
			s.Equal(tc.expected.ApprovalComment, result.ApprovalComment)
			s.Equal(tc.expected.PublishError, result.PublishError)

			// Verify expectations were met
			s.NoError(s.mock.ExpectationsWereMet())
		})
	}
}

// TestNewDraftRepository tests the repository constructor
func (s *DraftRepositorySuccessTestSuite) TestNewDraftRepository() {
	testCases := []struct {
		name        string
		db          *sql.DB
		expectError bool
	}{
		{
			name:        "Valid_database_connection_does_not_return_error",
			db:          s.db,
			expectError: false,
		},
		{
			name:        "Nil_database_connection_does_not_return_error",
			db:          nil,
			expectError: false, // The implementation doesn't check for nil DB
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			repo, err := repositories.NewDraftRepository(tc.db)

			// Assert
			if tc.expectError {
				s.Error(err)
				s.Nil(repo)
			} else {
				s.NoError(err)
				s.NotNil(repo)
			}
		})
	}
}

// TestDraftRepositorySuite runs all the test suites
func TestDraftRepositorySuite(t *testing.T) {
	suite.Run(t, new(DraftRepositorySuccessTestSuite))
	suite.Run(t, new(DraftRepositoryFailureTestSuite))
}
