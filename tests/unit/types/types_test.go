package types_test

import (
	"ISO_Auditing_Tool/pkg/types"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypesJSONMarshalUnmarshal(t *testing.T) {
	testCases := []struct {
		name string
		obj  interface{}
	}{
		{"Audit", createTestAudit()},
		{"ISOStandard", createTestISOStandard()},
		{"User", types.User{ID: "user_1", Name: "User 1"}},
		{"AuditQuestion", createTestAuditQuestion()},
		{"Evidence", types.Evidence{ID: 1, QuestionID: 1, Expected: "Expected Evidence"}},
		{"EvidenceProvided", types.EvidenceProvided{ID: 1, EvidenceID: 1, AuditQuestionID: 1, Provided: "Provided Evidence"}},
		{"Comment", types.Comment{ID: 1, UserID: "user_1", Text: "Comment 1", User: types.User{ID: "user_1", Name: "User 1"}}},
		{"Clause", createTestClause()},
		{"Section", createTestSection()},
		{"Subsection", createTestSubsection()},
		{"Question", types.Question{ID: 1, Text: "Question 1", Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testMarshalUnmarshal(t, tc.obj)
		})
	}
}

func testMarshalUnmarshal(t *testing.T, obj interface{}) {
	data, err := json.Marshal(obj)
	require.NoError(t, err, "Marshal should not return an error")

	unmarshalledObj := createNewInstance(obj)
	err = json.Unmarshal(data, unmarshalledObj)
	require.NoError(t, err, "Unmarshal should not return an error")

	expected := reflect.ValueOf(obj)
	actual := reflect.ValueOf(unmarshalledObj)

	if expected.Kind() == reflect.Ptr {
		expected = expected.Elem()
	}
	if actual.Kind() == reflect.Ptr {
		actual = actual.Elem()
	}

	assert.Equal(t, expected.Interface(), actual.Interface(), "Unmarshalled object should equal original object")
}

func createNewInstance(obj interface{}) interface{} {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Name() {
	case "Audit":
		return new(types.Audit)
	case "ISOStandard":
		return new(types.ISOStandard)
	case "User":
		return new(types.User)
	case "AuditQuestion":
		return new(types.AuditQuestion)
	case "Evidence":
		return new(types.Evidence)
	case "EvidenceProvided":
		return new(types.EvidenceProvided)
	case "Comment":
		return new(types.Comment)
	case "Clause":
		return new(types.Clause)
	case "Section":
		return new(types.Section)
	case "Subsection":
		return new(types.Subsection)
	case "Question":
		return new(types.Question)
	default:
		panic(fmt.Sprintf("Unknown type: %v", t.Name()))
	}
}

func createTestAudit() types.Audit {
	return types.Audit{
		ID:            1,
		Datetime:      time.Now().UTC(),
		ISOStandardID: 1,
		Name:          "Audit 1",
		Team:          "Team 1",
		UserID:        "user_1",
		ISOStandard:   types.ISOStandard{ID: 1, Name: "ISO 9001"},
		LeadAuditor:   types.User{ID: "user_1", Name: "Lead Auditor"},
		AuditQuestions: []types.AuditQuestion{
			createTestAuditQuestion(),
		},
	}
}

func createTestISOStandard() types.ISOStandard {
	return types.ISOStandard{
		ID:   1,
		Name: "ISO 9001",
		Clauses: []*types.Clause{
			createTestClause(),
		},
	}
}

func createTestAuditQuestion() types.AuditQuestion {
	return types.AuditQuestion{
		ID:         1,
		AuditID:    1,
		QuestionID: 1,
		EvidenceProvided: []types.EvidenceProvided{
			{ID: 1, EvidenceID: 1, AuditQuestionID: 1, Provided: "Evidence 1"},
		},
		Comments: []types.Comment{
			{ID: 1, UserID: "user_1", Text: "Comment 1", User: types.User{ID: "user_1", Name: "User 1"}},
		},
	}
}

func createTestClause() *types.Clause {
	return &types.Clause{
		ID:            1,
		ISOStandardID: 1,
		Name:          "Clause 1",
		Sections:      []*types.Section{createTestSection()},
	}
}

func createTestSection() *types.Section {
	return &types.Section{
		ID:        1,
		ClauseID:  1,
		Name:      "Section 1",
		Questions: []*types.Question{{ID: 1, Text: "Question 1", Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}}}},
	}
}

func createTestSubsection() types.Subsection {
	return types.Subsection{
		ID:        1,
		SectionID: 1,
		Name:      "Subsection 1",
		Questions: []*types.Question{{ID: 1, Text: "Question 1", Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}}}},
	}
}
