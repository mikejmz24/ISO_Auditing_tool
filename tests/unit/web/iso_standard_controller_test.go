package web_test

import (
	apiController "ISO_Auditing_Tool/cmd/api/controllers"
	webController "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/middleware"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/pkg/validators"
	"ISO_Auditing_Tool/tests/testutils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	// "log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// TestSuite combines all test suites
type TestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockRepo *testutils.MockIsoStandardRepository
	standard types.ISOStandard
	cleanup  []func()
}

type TestFormData struct {
	suite.Suite
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func loadTestData(filePath string) types.ISOStandard {
	data := getJSONData(filePath)

	var jsonData struct {
		ISOStandards []types.ISOStandard `json:"iso_standards"`
	}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal test data: %v", err))
	}

	if len(jsonData.ISOStandards) == 0 {
		panic("No ISO standards found in test data")
	}

	return jsonData.ISOStandards[0]
}

func getJSONData(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load data file: %v", err))
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("Failed to read data: %w", err))
	}
	return data
}

func (suite *TestSuite) SetupTest() {
	suite.cleanup = make([]func(), 0)
	testStandard := loadTestData("../../testdata/iso_standards_test01.json")
	suite.setupMockRepo()
	suite.setupRouter()
	suite.standard = testStandard

	validators.GetValidator()

	//	if err := validators.InitValidator(); err != nil {
	//		suite.T().Fatalf("Failed to initailize validator: %v", err)
	//	}
}

func (suite *TestSuite) setupMockRepo() {
	suite.mockRepo = new(testutils.MockIsoStandardRepository)
	if suite.mockRepo == nil {
		panic("mockRepo is nil")
	}
}

func (suite *TestSuite) setupRouter() {
	suite.router = setupRouter(suite.mockRepo)
	if suite.router == nil {
		panic("router is nil")
	}
}

func setupRouter(repo *testutils.MockIsoStandardRepository) *gin.Engine {
	apiController := apiController.NewApiIsoStandardController(repo)
	webController := webController.NewWebIsoStandardController(apiController)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	templatesPath := filepath.Join("..", "..", "..", "templates", "*.templ")
	router.LoadHTMLGlob(templatesPath)

	webGroup := router.Group("/web")
	{
		webGroup.GET("/iso_standards", webController.GetAllISOStandards)
		webGroup.GET("/iso_standards/:id", webController.GetISOStandardByID)
		webGroup.GET("/iso_standards/add", webController.RenderAddISOStandardForm)
		webGroup.POST("/iso_standards/add", webController.CreateISOStandard)
		webGroup.PUT("/iso_standards/:id", webController.UpdateISOStandard)
		webGroup.DELETE("/iso_standards/:id", webController.DeleteISOStandard)
	}
	return router
}

func (suite *TestSuite) TearDownTest() {
	suite.mockRepo.Reset()

	if len(suite.cleanup) > 0 {
		// Execute all cleanup functions in reverse order
		for i := len(suite.cleanup) - 1; 1 >= 0; i-- {
			if suite.cleanup[i] != nil {
				suite.cleanup[i]()
			}
		}
		// Clear the cleanup slice
		suite.cleanup = nil
	}

	// Clear the mock if it exists
	if suite.mockRepo != nil {
		suite.mockRepo.ExpectedCalls = nil
		suite.mockRepo.Calls = nil
	}

	// Clear the router
	suite.router = nil
}

func (suite *TestSuite) TestCreateISOStandard_Validation() {
	// Setup common test resources
	// validate := validator.New()
	// if err := types.RegisterCustomValidators(validate); err != nil {
	// 	log.Fatalf("Failed to register custom validators: %v", err)
	// }
	//
	validationTests := []struct {
		name          string
		setupRequest  func() (string, io.Reader, string)
		expectedError custom_errors.CustomError
	}{
		{
			name: "EmptyForm_ReturnsError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.EmptyData("Form"),
		},
		{
			name: "InvalidFormEncoding_ReturnsError",
			setupRequest: func() (string, io.Reader, string) {
				invalidBody := []byte("I'm not form-encoded!")
				return "/web/iso_standards/add", bytes.NewReader(invalidBody), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.ErrInvalidFormData,
		},
		{
			name: "MissingNameField_ReturnsError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"wrongField": {"ISO 9001"}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.EmptyField("string", "Name"),
		},
		{
			name: "EmptyNameField_ReturnsError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"name": {""}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.EmptyField("string", "Name"),
		},
		{
			name: "NameFieldWithOneCharacter_ReturnsMinError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"name": {"a"}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.MinFieldCharacters("Name", 2),
		},
		{
			name: "NameFieldWithOnehundredoneCharacters_ReturnsMaxError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"name": {"Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean ma"}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.MaxFieldCharacters("Name", 100),
		},
		{
			name: "NameFieldWithTrue_ReturnsIsABoolError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"name": {"true"}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.IsABool("Name")},
		{
			name: "NameFieldWithFalse_ReturnsIsABoolError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"name": {"false"}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.IsABool("Name"),
		},
		{
			name: "NameFieldWithOff_ReturnsIsABoolError",
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{"name": {"off"}}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: *custom_errors.IsABool("Name"),
		},
	}

	for _, tc := range validationTests {
		suite.Run(tc.name, func() {
			// Createa a fresh router for each test
			suite.SetupTest()

			// Setup
			path, body, contentType := tc.setupRequest()
			req, _ := http.NewRequest(http.MethodPost, path, body)
			req.Header.Set("Content-Type", contentType)

			// Create a new recorder for each test
			w := httptest.NewRecorder()

			// Execute the request
			suite.router.ServeHTTP(w, req)

			// Verify response
			expectedErr, _ := json.Marshal(tc.expectedError)

			fmt.Print("--- PRINT", string(expectedErr), "\n")
			fmt.Print("--- PRINT Body: ", w.Body.String(), "\n")

			suite.Equal(tc.expectedError.StatusCode, w.Code)
			suite.Equal(w.Body.String(), string(expectedErr))

			// Global cleanup after each test
			w.Flush()
			// suite.Contains(w.Body.String(), tc.expectedError.Message)
			// suite.Contains(w.Body.String(), tc.expectedError.Error())
		})
	}
}

// TODO : Remember to clear all resources to avoid 500 http server errors
func (suite *TestSuite) TestCreateISOStandard_Repository() {
	repositoryTests := []struct {
		name          string
		setupMock     func()
		setupRequest  func() (string, io.Reader, string)
		expectedError *custom_errors.CustomError
		validateExtra func(*httptest.ResponseRecorder)
	}{
		{
			name: "RepositoryError_ReturnsError",
			setupMock: func() {
				// Clear any existing mock expectations first
				suite.mockRepo.ExpectedCalls = nil
				suite.mockRepo.Calls = nil
				suite.mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(types.ISOStandard{}, errors.New("database error"))
			},
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{
					"name": {suite.standard.Name},
				}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			expectedError: custom_errors.NewCustomError(http.StatusInternalServerError, "Failed to create ISO standard", nil),
		},
		{
			name: "ValidData_CreatesAndRedirects",
			setupMock: func() {
				// Clear any esisting mock expectations first
				suite.mockRepo.ExpectedCalls = nil
				suite.mockRepo.Calls = nil

				standard := suite.standard
				suite.mockRepo.On("CreateISOStandard", mock.AnythingOfType("types.ISOStandard")).Return(standard, nil)
				suite.mockRepo.On("GetAllISOStandards").Return([]types.ISOStandard{standard}, nil)
			},
			setupRequest: func() (string, io.Reader, string) {
				formData := url.Values{
					"name": {suite.standard.Name},
				}
				return "/web/iso_standards/add", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded"
			},
			validateExtra: func(w *httptest.ResponseRecorder) {
				suite.Equal(http.StatusFound, w.Code)
				location := w.Header().Get("Location")
				suite.Equal("/web/iso_standards", location)

				req, err := http.NewRequest(http.MethodGet, location, nil)
				suite.NoError(err)

				resp := httptest.NewRecorder()
				suite.router.ServeHTTP(resp, req)

				suite.Equal(http.StatusOK, resp.Code)
				suite.Contains(resp.Body.String(), "ISO 9001")
			},
		},
	}

	for _, tc := range repositoryTests {
		suite.Run(tc.name, func() {
			// Reset mock and router for each test case
			suite.SetupTest()

			if tc.setupMock != nil {
				tc.setupMock()
			}

			path, body, contenType := tc.setupRequest()
			req, _ := http.NewRequest(http.MethodPost, path, body)
			req.Header.Set("Content-Type", contenType)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			if tc.expectedError != nil {
				suite.Equal(tc.expectedError.StatusCode, w.Code)
				suite.Contains(w.Body.String(), tc.expectedError.Error())
			}

			if tc.validateExtra != nil {
				tc.validateExtra(w)
			}

			suite.mockRepo.AssertExpectations(suite.T())
		})
	}
}

func TestWebISOStandardController(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func TestWebFormDataMethods(t *testing.T) {
	suite.Run(t, new(TestFormData))
}

func (suite *TestFormData) TestConversion() {
	type testCase[T any] struct {
		name   string
		input  T
		output url.Values
	}
	testCases := []testCase[any]{
		{
			name: "ISOStandardWithAllFields",
			input: types.ISOStandardForm{
				Name: "ISO TEST",
				Clauses: []*types.ClauseForm{
					{
						ID:            1,
						ISOStandardID: 1,
						Name:          "Clause Test",
						Sections: []*types.SectionForm{
							{
								ID:       1,
								ClauseID: 1,
								Name:     "Section Test",
								Questions: []*types.QuestionForm{
									{
										ID:           1,
										SectionID:    1,
										SubsectionID: 1,
										Text:         "Question Test",
										Evidence: []types.EvidenceForm{
											{
												ID:         1,
												QuestionID: 1,
												Expected:   "Evidence Test",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			output: func() url.Values {
				encoded := url.Values{}
				// Main ISO Standard fields
				// encoded.Set("id", "2")
				encoded.Set("name", "ISO TEST")

				// Clauses
				encoded.Set("clauses[0].id", "1")
				encoded.Set("clauses[0].iso_standard_id", "1")
				encoded.Set("clauses[0].name", "Clause Test")

				// Sections
				encoded.Set("clauses[0].sections[0].id", "1")
				encoded.Set("clauses[0].sections[0].clause_id", "1")
				encoded.Set("clauses[0].sections[0].name", "Section Test")

				// Questions
				encoded.Set("clauses[0].sections[0].questions[0].id", "1")
				encoded.Set("clauses[0].sections[0].questions[0].section_id", "1")
				encoded.Set("clauses[0].sections[0].questions[0].subsection_id", "1")
				encoded.Set("clauses[0].sections[0].questions[0].text", "Question Test")

				// Evidence
				encoded.Set("clauses[0].sections[0].questions[0].evidence[0].id", "1")
				encoded.Set("clauses[0].sections[0].questions[0].evidence[0].question_id", "1")
				encoded.Set("clauses[0].sections[0].questions[0].evidence[0].expected", "Evidence Test")

				return encoded
			}(),
		},
		{
			name:  "ISOStandardWithIDAndName",
			input: types.ISOStandardForm{Name: "ISO TEST"},
			output: func() url.Values {
				encoded := url.Values{}
				encoded.Set("name", "ISO TEST")
				return encoded
			}(),
		},
		{
			name:  "ISOStandardWithOnlyName",
			input: types.ISOStandardForm{Name: "ISO TEST"},
			output: func() url.Values {
				encoded := url.Values{}
				encoded.Set("name", "ISO TEST")
				return encoded
			}(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			formData, _ := webController.ConvertStructToForm(tc.input)
			suite.Equal(tc.output, formData)
		})
	}
}
