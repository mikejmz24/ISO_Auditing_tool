package testutils

import (
	"ISO_Auditing_Tool/pkg/types"
	"encoding/json"
	"time"
)

// CreateTestDraft creates a test draft for use in tests
func CreateTestDraft() types.Draft {
	now := time.Now().UTC().Truncate(time.Second)
	return types.Draft{
		ID:              1,
		TypeID:          1,
		ObjectID:        42,
		StatusID:        1,
		Version:         1,
		Data:            json.RawMessage(`{"name": "ISO 27001", "description": "Information Security Standard"}`),
		Diff:            json.RawMessage(`{"name": {"old": "ISO 27000", "new": "ISO 27001"}}`),
		UserID:          10,
		ApproverID:      0,
		ApprovalComment: "",
		PublishError:    "",
		CreatedAt:       now,
		UpdatedAt:       now,
		ExpiresAt:       now.Add(7 * 24 * time.Hour),
	}
}

// CreateTestStandard creates a test ISO standard
func CreateTestStandard() types.Standard {
	return types.Standard{
		ID:   1,
		Name: "ISO 9001",
		Requirements: []types.Requirement{
			{
				ID:         1,
				StandardID: 1,
				Name:       "Quality Management System",
			},
		},
	}
}
