package types_test

import (
	"ISO_Auditing_Tool/pkg/types"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestTypesSuite is the main suite container
type TestTypesSuite struct {
	suite.Suite
}

// TestTypesJSONMarshalling handles JSON marshal tests
type TestTypesJSONMarshalling struct {
	suite.Suite
}

// TestTypesJSONUnmarshalling handles JSON unmarshal tests
type TestTypesJSONUnmarshalling struct {
	suite.Suite
}

// TestTypesFormEncoding handles form encoding/decoding tests
type TestTypesFormEncoding struct {
	suite.Suite
}

// Test cases struct to be shared across suites
type testCase struct {
	name string
	obj  interface{}
}

// TestTypesToFormMethods
type TestTypesToFormMethods struct {
	suite.Suite
}

// Generic standalone helper function
func runDynamicTypeConversion[TInput any, TOutput any](
	t *testing.T,
	cases []struct {
		name         string
		input        TInput
		expectedFunc func(input TInput) TOutput
		expectedData TOutput
	},
) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the conversion function
			result := tc.expectedFunc(tc.input)

			// Assertions
			require.NotNil(t, result, "Result should not be nil")
			assert.True(t, reflect.DeepEqual(tc.expectedData, result), "Expected: %+v, Got: %+v", tc.expectedData, result)
		})
	}
}

func (suite *TestTypesToFormMethods) TestISOStandardFormConversion() {
	// Define test cases
	isoStandardCases := []struct {
		name         string
		input        types.ISOStandardForm
		expectedFunc func(input types.ISOStandardForm) *types.ISOStandard
		expectedData *types.ISOStandard
	}{
		{
			name: "ISOStandardForm",
			input: types.ISOStandardForm{
				Name: &[]string{"ISO 9001"}[0],
			},
			expectedFunc: func(input types.ISOStandardForm) *types.ISOStandard {
				return input.ToISOStandard()
			},
			expectedData: &types.ISOStandard{
				Name: "ISO 9001",
			},
		},
	}

	runDynamicTypeConversion(suite.T(), isoStandardCases)
}

type FormEncoder interface {
	EncodeForm() url.Values
}

// Example implementation for ISOStandardForm
func (f *TestTypesFormEncoding) EncodeForm() url.Values {
	values := url.Values{}

	v := reflect.ValueOf(f).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := t.Field(i)

		// Dynamic tag handling
		tag := typeField.Tag.Get("form")
		if tag == "" {
			tag = typeField.Name
		}

		// Nil-safe encoding
		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				values.Set(tag, fmt.Sprintf("%v", field.Elem().Interface()))
			}
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Ptr {
					if !elem.IsNil() {
						values.Add(tag, fmt.Sprintf("%v", elem.Elem().Interface()))
					}
				} else {
					values.Add(tag, fmt.Sprintf("%v", elem.Interface()))
				}
			}
		}
	}

	return values
}

// getTestCases returns common test cases for all suites
func getTestCases() []testCase {
	return []testCase{
		{name: "Audit", obj: createTestAudit()},
		{name: "ISOStandard", obj: createTestISOStandard()},
		{name: "ISOStandardWithJustName", obj: createTestISOStandardWithOnlyName()},
		{name: "User", obj: types.User{ID: "user_1", Name: "User 1"}},
		{name: "AuditQuestion", obj: createTestAuditQuestion()},
		{name: "Evidence", obj: types.Evidence{ID: 1, QuestionID: 1, Expected: "Expected Evidence"}},
		{name: "EvidenceProvided", obj: types.EvidenceProvided{ID: 1, EvidenceID: 1, AuditQuestionID: 1, Provided: "Provided Evidence"}},
		{name: "Comment", obj: types.Comment{ID: 1, UserID: "user_1", Text: "Comment 1", User: types.User{ID: "user_1", Name: "User 1"}}},
		{name: "Clause", obj: createTestClause()},
		{name: "Section", obj: createTestSection()},
		{name: "Subsection", obj: createTestSubsection()},
		{name: "Question", obj: types.Question{ID: 1, Text: "Question 1", Evidence: []types.Evidence{{ID: 1, QuestionID: 1, Expected: "Evidence 1"}}}},
	}
}

// TestMarshalling tests JSON marshalling
func (suite *TestTypesJSONMarshalling) TestMarshalling() {
	for _, tc := range getTestCases() {
		suite.Run(tc.name, func() {
			data, err := json.Marshal(tc.obj)
			require.NoError(suite.T(), err, "Marshal should not return an error")
			assert.NotEmpty(suite.T(), data, "Marshalled data should not be empty")
		})
	}
}

// TestUnmarshalling tests JSON unmarshalling
func (suite *TestTypesJSONUnmarshalling) TestUnmarshalling() {
	for _, tc := range getTestCases() {
		suite.Run(tc.name, func() {
			data, err := json.Marshal(tc.obj)
			require.NoError(suite.T(), err, "Marshal should not return an error")

			unmarshalledObj := createNewInstance(tc.obj)
			err = json.Unmarshal(data, unmarshalledObj)
			require.NoError(suite.T(), err, "Unmarshal should not return an error")

			expected := reflect.ValueOf(tc.obj)
			actual := reflect.ValueOf(unmarshalledObj)

			if expected.Kind() == reflect.Ptr {
				expected = expected.Elem()
			}
			if actual.Kind() == reflect.Ptr {
				actual = actual.Elem()
			}

			assert.Equal(suite.T(), expected.Interface(), actual.Interface(), "Unmarshalled object should equal original object")
		})
	}
}

// TestTypes runs all the suites
func TestTypes(t *testing.T) {
	suite.Run(t, new(TestTypesJSONMarshalling))
	suite.Run(t, new(TestTypesJSONUnmarshalling))
	suite.Run(t, new(TestTypesFormEncoding))
	suite.Run(t, new(TestTypesToFormMethods))
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

func createTestISOStandardWithOnlyName() types.ISOStandard {
	return types.ISOStandard{
		Name: "ISO 9001",
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
