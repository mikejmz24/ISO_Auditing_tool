package types_test

import (
	"ISO_Auditing_Tool/pkg/types"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TypesTestSuite struct {
	suite.Suite
}

// Helper function to compare time.Time fields
func timesAreEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

func (suite *TypesTestSuite) TestAudit_MarshalUnmarshal() {
	audit := types.Audit{
		ID:            1,
		Datetime:      time.Now().UTC(),
		ISOStandardID: 1,
		Name:          "Audit 1",
		Team:          "Team 1",
		UserID:        "user_1",
		ISOStandard: types.ISOStandard{
			ID:   1,
			Name: "ISO 9001",
		},
		LeadAuditor: types.User{
			ID:   "user_1",
			Name: "Lead Auditor",
		},
		AuditQuestions: []types.AuditQuestion{
			{
				ID:         1,
				AuditID:    1,
				QuestionID: 1,
				EvidenceProvided: []types.EvidenceProvided{
					{
						ID:              1,
						EvidenceID:      1,
						AuditQuestionID: 1,
						Provided:        "Evidence 1",
					},
				},
				Comments: []types.Comment{
					{
						ID:     1,
						UserID: "user_1",
						Text:   "Comment 1",
						User: types.User{
							ID:   "user_1",
							Name: "User 1",
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(audit)
	assert.NoError(suite.T(), err)

	var unmarshalledAudit types.Audit
	err = json.Unmarshal(data, &unmarshalledAudit)
	assert.NoError(suite.T(), err)

	// Compare fields manually
	assert.Equal(suite.T(), audit.ID, unmarshalledAudit.ID)
	assert.True(suite.T(), timesAreEqual(audit.Datetime, unmarshalledAudit.Datetime))
	assert.Equal(suite.T(), audit.ISOStandardID, unmarshalledAudit.ISOStandardID)
	assert.Equal(suite.T(), audit.Name, unmarshalledAudit.Name)
	assert.Equal(suite.T(), audit.Team, unmarshalledAudit.Team)
	assert.Equal(suite.T(), audit.UserID, unmarshalledAudit.UserID)
	assert.Equal(suite.T(), audit.ISOStandard, unmarshalledAudit.ISOStandard)
	assert.Equal(suite.T(), audit.LeadAuditor, unmarshalledAudit.LeadAuditor)
	assert.Equal(suite.T(), audit.AuditQuestions, unmarshalledAudit.AuditQuestions)
}

func (suite *TypesTestSuite) TestISOStandard_MarshalUnmarshal() {
	isoStandard := types.ISOStandard{
		ID:   1,
		Name: "ISO 9001",
		Clauses: []*types.Clause{
			{
				ID:            1,
				ISOStandardID: 1,
				Name:          "Clause 1",
				Sections: []*types.Section{
					{
						ID:       1,
						ClauseID: 1,
						Name:     "Section 1",
						Questions: []*types.Question{
							{
								ID:       1,
								Text:     "Question 1",
								Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}},
							},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(isoStandard)
	assert.NoError(suite.T(), err)

	var unmarshalledISOStandard types.ISOStandard
	err = json.Unmarshal(data, &unmarshalledISOStandard)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), isoStandard, unmarshalledISOStandard)
}

func (suite *TypesTestSuite) TestUser_MarshalUnmarshal() {
	user := types.User{
		ID:   "user_1",
		Name: "User 1",
	}

	data, err := json.Marshal(user)
	assert.NoError(suite.T(), err)

	var unmarshalledUser types.User
	err = json.Unmarshal(data, &unmarshalledUser)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), user, unmarshalledUser)
}

func (suite *TypesTestSuite) TestAuditQuestion_MarshalUnmarshal() {
	auditQuestion := types.AuditQuestion{
		ID:         1,
		AuditID:    1,
		QuestionID: 1,
		EvidenceProvided: []types.EvidenceProvided{
			{
				ID:              1,
				EvidenceID:      1,
				AuditQuestionID: 1,
				Provided:        "Evidence 1",
			},
		},
		Comments: []types.Comment{
			{
				ID:     1,
				UserID: "user_1",
				Text:   "Comment 1",
				User: types.User{
					ID:   "user_1",
					Name: "User 1",
				},
			},
		},
	}

	data, err := json.Marshal(auditQuestion)
	assert.NoError(suite.T(), err)

	var unmarshalledAuditQuestion types.AuditQuestion
	err = json.Unmarshal(data, &unmarshalledAuditQuestion)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), auditQuestion, unmarshalledAuditQuestion)
}

func (suite *TypesTestSuite) TestEvidence_MarshalUnmarshal() {
	evidence := types.Evidence{
		ID:         1,
		QuestionID: 1,
		Expected:   "Expected Evidence",
	}

	data, err := json.Marshal(evidence)
	assert.NoError(suite.T(), err)

	var unmarshalledEvidence types.Evidence
	err = json.Unmarshal(data, &unmarshalledEvidence)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), evidence, unmarshalledEvidence)
}

func (suite *TypesTestSuite) TestEvidenceProvided_MarshalUnmarshal() {
	evidenceProvided := types.EvidenceProvided{
		ID:              1,
		EvidenceID:      1,
		AuditQuestionID: 1,
		Provided:        "Provided Evidence",
	}

	data, err := json.Marshal(evidenceProvided)
	assert.NoError(suite.T(), err)

	var unmarshalledEvidenceProvided types.EvidenceProvided
	err = json.Unmarshal(data, &unmarshalledEvidenceProvided)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), evidenceProvided, unmarshalledEvidenceProvided)
}

func (suite *TypesTestSuite) TestComment_MarshalUnmarshal() {
	comment := types.Comment{
		ID:     1,
		UserID: "user_1",
		Text:   "Comment 1",
		User: types.User{
			ID:   "user_1",
			Name: "User 1",
		},
	}

	data, err := json.Marshal(comment)
	assert.NoError(suite.T(), err)

	var unmarshalledComment types.Comment
	err = json.Unmarshal(data, &unmarshalledComment)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), comment, unmarshalledComment)
}

func (suite *TypesTestSuite) TestClause_MarshalUnmarshal() {
	clause := types.Clause{
		ID:            1,
		ISOStandardID: 1,
		Name:          "Clause 1",
		Sections: []*types.Section{
			{
				ID:       1,
				ClauseID: 1,
				Name:     "Section 1",
				Questions: []*types.Question{
					{
						ID:       1,
						Text:     "Question 1",
						Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}},
					},
				},
			},
		},
	}

	data, err := json.Marshal(clause)
	assert.NoError(suite.T(), err)

	var unmarshalledClause types.Clause
	err = json.Unmarshal(data, &unmarshalledClause)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), clause, unmarshalledClause)
}

func (suite *TypesTestSuite) TestSection_MarshalUnmarshal() {
	section := types.Section{
		ID:       1,
		ClauseID: 1,
		Name:     "Section 1",
		Questions: []*types.Question{
			{
				ID:       1,
				Text:     "Question 1",
				Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}},
			},
		},
	}

	data, err := json.Marshal(section)
	assert.NoError(suite.T(), err)

	var unmarshalledSection types.Section
	err = json.Unmarshal(data, &unmarshalledSection)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), section, unmarshalledSection)
}

func (suite *TypesTestSuite) TestSubsection_MarshalUnmarshal() {
	subsection := types.Subsection{
		ID:        1,
		SectionID: 1,
		Name:      "Subsection 1",
		Questions: []*types.Question{
			{
				ID:       1,
				Text:     "Question 1",
				Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}},
			},
		},
	}

	data, err := json.Marshal(subsection)
	assert.NoError(suite.T(), err)

	var unmarshalledSubsection types.Subsection
	err = json.Unmarshal(data, &unmarshalledSubsection)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), subsection, unmarshalledSubsection)
}

func (suite *TypesTestSuite) TestQuestion_MarshalUnmarshal() {
	question := types.Question{
		ID:       1,
		Text:     "Question 1",
		Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}},
	}

	data, err := json.Marshal(question)
	assert.NoError(suite.T(), err)

	var unmarshalledQuestion types.Question
	err = json.Unmarshal(data, &unmarshalledQuestion)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), question, unmarshalledQuestion)
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}
