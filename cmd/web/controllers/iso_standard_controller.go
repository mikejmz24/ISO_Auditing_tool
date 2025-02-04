package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	// "strconv"
	"strings"

	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/types"
	"ISO_Auditing_Tool/pkg/validators"
	"ISO_Auditing_Tool/templates"
	"ISO_Auditing_Tool/templates/iso_standards"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// Interface for the API controller to allow for easier testing and mocking
type ApiIsoStandardController interface {
	GetAllISOStandards(c *gin.Context)
	GetISOStandardByID(c *gin.Context)
	CreateISOStandard(c *gin.Context)
	UpdateISOStandard(c *gin.Context)
	DeleteISOStandard(c *gin.Context)
}

type WebIsoStandardController struct {
	ApiController ApiIsoStandardController
}

func NewWebIsoStandardController(apiController ApiIsoStandardController) *WebIsoStandardController {
	return &WebIsoStandardController{ApiController: apiController}
}

func (wc *WebIsoStandardController) GetAllISOStandards(c *gin.Context) {
	isoStandards, err := wc.fetchAllISOStandards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ISO standards"})
		return
	}
	templ.Handler(templates.Base(iso_standards.List(isoStandards))).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) GetISOStandardByID(c *gin.Context) {
	id := c.Param("id")
	isoStandard, err := wc.fetchISOStandardByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ISO standard not found"})
		return
	}
	templ.Handler(templates.Base(iso_standards.Detail(*isoStandard))).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) RenderAddISOStandardForm(c *gin.Context) {
	templ.Handler(templates.Base(iso_standards.Add())).ServeHTTP(c.Writer, c.Request)
}

func (wc *WebIsoStandardController) CreateISOStandard(c *gin.Context) {
	// Read the raw body first
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.ErrInvalidFormData)
		return
	}

	// Restore the body for later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	if len(rawBody) == 0 {
		err := custom_errors.EmptyData("Form")
		c.JSON(err.StatusCode, err)
		return
	}

	// Check if the body contains an equals sign (=) which is required for form encoding
	if !bytes.Contains(rawBody, []byte("=")) {
		c.JSON(http.StatusBadRequest, custom_errors.ErrInvalidFormData)
		return
	}

	// Try to parse as form values
	_, err = url.ParseQuery(string(rawBody))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.ErrInvalidFormData)
		return
	}
	var formData types.ISOStandardForm
	// Parse and bind form data
	if err := c.Bind(&formData); err != nil {
		log.Printf("Error binding form data: %v", err)
		c.JSON(http.StatusBadRequest, custom_errors.ErrInvalidFormData)
		return
	}

	// Validate form data
	// if err := validateFormData(c, &formData); err != nil {
	if err := validateFormData(&formData); err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	// Convert to JSON and forward to API
	isoStandard := formData.ToISOStandard()
	jsonData, err := json.Marshal(isoStandard)
	if err != nil {
		log.Printf("Error marshalling form data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process form data"})
		return
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/iso_standards", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = req

	wc.ApiController.CreateISOStandard(apiContext)

	if recorder.Code == http.StatusCreated || recorder.Code == http.StatusOK {
		c.Redirect(http.StatusFound, "/web/iso_standards")
	} else {
		c.JSON(recorder.Code, gin.H{"error": "Failed to create ISO standard"})
	}
}

func validateFormData(formData interface{}) *custom_errors.CustomError {
	return validators.ValidateStruct(formData)
}

// Separate validation function for better organization and reusability
// func validateFormData(c *gin.Context, formData *types.ISOStandardForm) *custom_errors.CustomError {
// 	if _, exists := c.Request.PostForm["name"]; !exists {
// 		return custom_errors.MissingField("name")
// 	}
//
// 	if formData == nil {
// 		return custom_errors.EmptyData("Form")
// 	}
//
// 	if formData.Name == "" {
// 		return custom_errors.EmptyField("string", "name")
// 	}
//
// 	if _, exists := c.Request.PostForm["name"]; !exists {
// 		// return custom_errors.MissingField("name")
// 		return custom_errors.EmptyField("string", "name")
// 	}
//
// 	if isInvalidString(formData.Name) {
// 		return custom_errors.InvalidDataType("name", "string")
// 	}
//
// 	return nil

// if err := c.Bind(formData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Form was not binded :("})
// 	}
//
// 	// Empty string
// 	if formData.Name == "" {
// 		return custom_errors.EmptyField("string", "name")
// 	}
// 	validate := validator.New()
// 	if err := validate.Struct(formData); err != nil {
// 		// Process validation errors
// 		if validationErrors, ok := err.(validator.ValidationErrors); ok {
// 			for _, e := range validationErrors {
// 				switch e.Tag() {
// 				case "required":
// 					return custom_errors.EmptyField("string", "name")
// 				case "min":
// 					return custom_errors.MinFieldCharacters("name", 2)
// 				case "max":
// 					return custom_errors.MaxFieldCharacters("name", 100)
// 				case "not_boolean":
// 					return custom_errors.NewCustomError(201, "I was not a boolean", nil)
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

func (wc *WebIsoStandardController) UpdateISOStandard(c *gin.Context) {
	var formData map[string]string
	if err := c.Bind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	jsonData, err := json.Marshal(formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal form data"})
		return
	}

	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer
	apiContext.Request.Body = io.NopCloser(strings.NewReader(string(jsonData)))
	apiContext.Request.Header.Set("Content-Type", "application/json")

	wc.ApiController.UpdateISOStandard(apiContext)

	if apiContext.Writer.Status() == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	} else {
		c.JSON(apiContext.Writer.Status(), gin.H{"error": "Failed to update ISO standard"})
	}
}

func (wc *WebIsoStandardController) DeleteISOStandard(c *gin.Context) {
	apiContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	apiContext.Request = c.Request
	apiContext.Writer = c.Writer

	wc.ApiController.DeleteISOStandard(apiContext)

	if apiContext.Writer.Status() == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	} else {
		c.JSON(apiContext.Writer.Status(), gin.H{"error": "Failed to delete ISO standard"})
	}
}

// Helper functions for fetching data from the API controller
func (wc *WebIsoStandardController) fetchAllISOStandards() ([]types.ISOStandard, error) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/iso_standards", nil)
	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = req

	wc.ApiController.GetAllISOStandards(apiContext)

	if recorder.Code != http.StatusOK {
		log.Printf("Error fetching ISO standards: %s", recorder.Body.String())
		return nil, fmt.Errorf("error fetching ISO standards")
	}

	var isoStandards []types.ISOStandard
	if err := json.Unmarshal(recorder.Body.Bytes(), &isoStandards); err != nil {
		log.Printf("Error unmarshalling ISO standards: %v", err)
		return nil, err
	}
	return isoStandards, nil
}

// Helper function to fetch a single ISO standard
func (wc *WebIsoStandardController) fetchISOStandardByID(id string) (*types.ISOStandard, error) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/iso_standards/"+id, nil)
	apiContext, _ := gin.CreateTestContext(recorder)
	apiContext.Request = req
	apiContext.Params = gin.Params{{Key: "id", Value: id}}

	wc.ApiController.GetISOStandardByID(apiContext)

	if recorder.Code != http.StatusOK {
		log.Printf("Error fetching ISO standard by ID: %s", recorder.Body.String())
		return nil, fmt.Errorf("error fetching ISO standard by ID")
	}

	var isoStandard types.ISOStandard
	if err := json.Unmarshal(recorder.Body.Bytes(), &isoStandard); err != nil {
		log.Printf("Error unmarshalling ISO standard: %v", err)
		return nil, err
	}
	return &isoStandard, nil
}

// func isInvalidString(input string) bool {
// 	if input == "true" || input == "false" {
// 		return true
// 	}
//
// 	if _, err := strconv.Atoi(input); err == nil {
// 		return true
// 	}
//
// 	if _, err := strconv.ParseFloat(input, 64); err == nil {
// 		return true
// 	}
//
// 	return false
// }

// ConvertStructToForm converts a struct into url.Values
func ConvertStructToForm(data interface{}) (url.Values, error) {
	values := url.Values{}
	if err := convertField(values, "", reflect.ValueOf(data)); err != nil {
		return nil, err
	}
	return values, nil
}

func convertField(values url.Values, prefix string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	t := v.Type()

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			tag := field.Tag.Get("form")
			if tag == "" {
				continue
			}

			tagParts := strings.Split(tag, ",")
			formKey := tagParts[0]

			// Construct the full key with prefix
			fullKey := formKey
			if prefix != "" {
				fullKey = prefix + "." + formKey
			}

			// Special handling for top-level ID fields
			if formKey == "id" && prefix == "" {
				if fieldValue.Kind() == reflect.Int || fieldValue.Kind() == reflect.Int64 {
					if fieldValue.Int() == 0 {
						continue
					} else if fieldValue.Int() == 1 {
						values.Set(formKey, "")
						continue
					}
				}
			}

			if err := convertField(values, fullKey, fieldValue); err != nil {
				return err
			}
		}

	case reflect.Slice:
		if v.IsNil() || v.Len() == 0 {
			return nil
		}

		for i := 0; i < v.Len(); i++ {
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			if err := convertField(values, newPrefix, v.Index(i)); err != nil {
				return err
			}
		}

	default:
		if !v.IsValid() {
			return nil
		}

		var value string
		switch v.Kind() {
		case reflect.String:
			value = v.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = fmt.Sprintf("%v", v.Int())
		case reflect.Float32, reflect.Float64:
			value = fmt.Sprintf("%g", v.Float())
		default:
			if stringer, ok := v.Interface().(fmt.Stringer); ok {
				value = stringer.String()
			} else {
				value = fmt.Sprintf("%v", v.Interface())
			}
		}

		values.Set(prefix, value)
	}

	return nil
}
