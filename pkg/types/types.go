package types

import (
	// "ISO_Auditing_Tool/pkg/custom_errors"
	"encoding/json"
	"errors"
	"fmt"
	// "net/http"
	"net/url"
	"reflect"
	"strconv"
	// "strings"
	"time"
	// "github.com/go-playground/validator/v10"
)

// FormValidator interface for form validation
type FormValidator interface {
	Validate() error
}

// Common form validation errors
var (
	ErrRequired = errors.New("field is required")
	ErrInvalid  = errors.New("invalid value")
)

// DecodeForm decodes url.Values into a struct using form tags
func DecodeForm(values url.Values, dst interface{}) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}
	v = v.Elem()

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		value := values.Get(formTag)
		if value == "" {
			continue
		}

		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int, reflect.Int64:
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid int value for field %s: %w", field.Name, err)
			}
			fieldValue.SetInt(intVal)
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid float value for field %s: %w", field.Name, err)
			}
			fieldValue.SetFloat(floatVal)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid bool value for field %s: %w", field.Name, err)
			}
			fieldValue.SetBool(boolVal)
		}
	}

	return nil
}

type CommentForm struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
	Text   string `json:"text" form:"text" binding:"required"`
}

func (f *CommentForm) Validate() error {
	if f.UserID == "" {
		return fmt.Errorf("user_id: %w", ErrRequired)
	}
	if f.Text == "" {
		return fmt.Errorf("text: %w", ErrRequired)
	}
	return nil
}

func UnmarshalAudit(data []byte) (Audit, error) {
	var r Audit
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Audit) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Audit struct {
	ID             int             `json:"id"`
	Datetime       time.Time       `json:"datetime"`
	ISOStandardID  int             `json:"iso_standard_id"`
	Name           string          `json:"name"`
	Team           string          `json:"team"`
	UserID         string          `json:"user_id"`
	ISOStandard    ISOStandard     `json:"iso_standard"`
	LeadAuditor    User            `json:"user"`
	AuditQuestions []AuditQuestion `json:"audit_questions"`
}

type AuditQuestion struct {
	ID               int                `json:"id"`
	AuditID          int                `json:"audit_id"`
	QuestionID       int                `json:"question_id"`
	EvidenceProvided []EvidenceProvided `json:"evidence_provided"`
	Comments         []Comment          `json:"comments"`
}

type ISOStandard struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Clauses []*Clause `json:"clauses,omitempty"`
}

type ISOStandardForm struct {
	// Name    string        `form:"name" validate:"required,min=3,max=100,not_boolean"`
	Name    string        `form:"name" validate:"required,min=3,max=100,not_boolean"`
	Clauses []*ClauseForm `form:"clauses,omitempty"`
}

type Clause struct {
	ID            int        `json:"id"`
	ISOStandardID int        `json:"iso_standard_id"`
	Name          string     `json:"name"`
	Sections      []*Section `json:"sections,omitempty"`
}

type ClauseForm struct {
	ID            int            `json:"id" form:"id"`
	ISOStandardID int            `json:"iso_standard_id" form:"iso_standard_id"`
	Name          string         `json:"name" form:"name" binding:"required"`
	Sections      []*SectionForm `json:"sections" form:"sections"`
}

type Section struct {
	ID        int         `json:"id"`
	ClauseID  int         `json:"clause_id"`
	Name      string      `json:"name"`
	Questions []*Question `json:"questions,omitempty"`
}

type SectionForm struct {
	ID        int             `json:"id" form:"id"`
	ClauseID  int             `json:"clause_id" form:"clause_id"`
	Name      string          `json:"name" form:"name" binding:"required"`
	Questions []*QuestionForm `json:"questions" form:"questions"`
}

type Subsection struct {
	ID        int         `json:"id"`
	SectionID int         `json:"section_id"`
	Name      string      `json:"name"`
	Questions []*Question `json:"questions,omitempty"`
}

type Question struct {
	ID           int        `json:"id"`
	SectionID    int        `json:"section_id,omitempty"`
	SubsectionID int        `json:"subsection_id,omitempty"`
	Text         string     `json:"text"`
	Evidence     []Evidence `json:"evidence,omitempty"`
}

type QuestionForm struct {
	ID           int            `form:"id"`
	SectionID    int            `form:"section_id,omitempty"`
	SubsectionID int            `form:"subsection_id,omitempty"`
	Text         string         `form:"text"`
	Evidence     []EvidenceForm `form:"evidence,omitempty"`
}

type Evidence struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Expected   string `json:"expected"`
}

type EvidenceForm struct {
	ID         int    `form:"id"`
	QuestionID int    `form:"question_id"`
	Expected   string `form:"expected"`
}

type EvidenceProvided struct {
	ID              int    `json:"id"`
	EvidenceID      int    `json:"evidence_id"`
	AuditQuestionID int    `json:"audit_question_id"`
	Provided        string `json:"provided"`
}

type Comment struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	Text   string `json:"text"`
	User   User   `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (f *ISOStandardForm) ToISOStandard() *ISOStandard {
	var clauses []*Clause
	for _, clauseForm := range f.Clauses {
		clauses = append(clauses, clauseForm.ToClause())
	}

	return &ISOStandard{
		Name:    f.Name,
		Clauses: clauses,
	}
}

func (f *ISOStandardForm) FromISOStandard(iso *ISOStandard) {
	f.Name = iso.Name

	var clauses []*ClauseForm
	for _, clause := range iso.Clauses {
		clauses = append(clauses, clause.ToClauseForm())
	}
	f.Clauses = clauses
}

func (c *Clause) ToClauseForm() *ClauseForm {
	var items []*SectionForm
	for _, item := range c.Sections {
		items = append(items, item.ToSectionForm())
	}

	return &ClauseForm{
		ID:            c.ID,
		ISOStandardID: c.ISOStandardID,
		Name:          c.Name,
		Sections:      items,
	}
}

func (c *Clause) FromClauseForm(iso *ClauseForm) {
	c.ID = iso.ID
	c.ISOStandardID = iso.ISOStandardID
	c.Name = iso.Name

	var items []*Section
	for _, item := range iso.Sections {
		items = append(items, item.ToSection())
	}
	c.Sections = items
}

func (c *ClauseForm) ToClause() *Clause {
	var items []*Section
	for _, item := range c.Sections {
		items = append(items, item.ToSection())
	}

	return &Clause{
		ID:            c.ID,
		ISOStandardID: c.ISOStandardID,
		Name:          c.Name,
		Sections:      items,
	}
}

func (f *ClauseForm) FromClause(iso *Clause) {
	f.ID = iso.ID
	f.ISOStandardID = iso.ISOStandardID
	f.Name = iso.Name

	var items []*SectionForm
	for _, item := range iso.Sections {
		items = append(items, item.ToSectionForm())
	}
	f.Sections = items
}

func (s *Section) ToSectionForm() *SectionForm {
	var items []*QuestionForm
	for _, item := range s.Questions {
		items = append(items, item.ToQuestionForm())
	}

	return &SectionForm{
		ID:        s.ID,
		ClauseID:  s.ClauseID,
		Name:      s.Name,
		Questions: items,
	}
}

func (c *Section) FromSectionForm(iso *SectionForm) {
	c.ID = iso.ID
	c.ClauseID = iso.ClauseID
	c.Name = iso.Name

	var items []*Question
	for _, item := range iso.Questions {
		items = append(items, item.ToQuestion())
	}
	c.Questions = items
}

func (c *SectionForm) ToSection() *Section {
	var items []*Question
	for _, item := range c.Questions {
		items = append(items, item.ToQuestion())
	}

	return &Section{
		ID:        c.ID,
		ClauseID:  c.ClauseID,
		Name:      c.Name,
		Questions: items,
	}
}

func (f *SectionForm) FromSection(iso *Section) {
	f.ID = iso.ID
	f.ClauseID = iso.ClauseID
	f.Name = iso.Name

	var items []*QuestionForm
	for _, item := range iso.Questions {
		items = append(items, item.ToQuestionForm())
	}
	f.Questions = items
}

func (i *Question) ToQuestionForm() *QuestionForm {
	var items []EvidenceForm
	for _, item := range i.Evidence {
		items = append(items, item.ToEvidenceForm())
	}

	return &QuestionForm{
		ID:           i.ID,
		SectionID:    i.SectionID,
		SubsectionID: i.SubsectionID,
		Text:         i.Text,
		Evidence:     items,
	}
}

func (i *Question) FromQuestionForm(iso *QuestionForm) {
	i.ID = iso.ID
	i.SectionID = iso.SectionID
	i.SubsectionID = iso.SubsectionID
	i.Text = iso.Text

	var items []Evidence
	for _, item := range iso.Evidence {
		items = append(items, item.ToEvidence())
	}
	i.Evidence = items
}

func (i *QuestionForm) ToQuestion() *Question {
	var items []Evidence
	for _, item := range i.Evidence {
		items = append(items, item.ToEvidence())
	}

	return &Question{
		ID:           i.ID,
		SectionID:    i.SectionID,
		SubsectionID: i.SubsectionID,
		Text:         i.Text,
		Evidence:     items,
	}
}

func (i *QuestionForm) FromQuestion(iso *Question) {
	i.ID = iso.ID
	i.SectionID = iso.SectionID
	i.SubsectionID = iso.SubsectionID
	i.Text = iso.Text

	var items []EvidenceForm
	for _, item := range iso.Evidence {
		items = append(items, item.ToEvidenceForm())
	}
	i.Evidence = items
}

func (i *Evidence) ToEvidenceForm() EvidenceForm {
	return EvidenceForm{
		ID:         i.ID,
		QuestionID: i.QuestionID,
		Expected:   i.Expected,
	}
}

func (i *Evidence) FromEvidenceForm(iso EvidenceForm) {
	i.ID = iso.ID
	i.QuestionID = iso.QuestionID
	i.Expected = iso.Expected
}

func (i *EvidenceForm) ToEvidence() Evidence {
	return Evidence{
		ID:         i.ID,
		QuestionID: i.QuestionID,
		Expected:   i.Expected,
	}
}

func (i *EvidenceForm) FromEvidence(iso Evidence) {
	i.ID = iso.ID
	i.QuestionID = iso.QuestionID
	i.Expected = iso.Expected
}

// Create a package-level validator instance
// var validate *validator.Validate

// Initialize validator withi custom validations
// func InitValidator() error {
// 	validate = validator.New()
//
// 	if err := RegisterCustomValidators(validate); err != nil {
// 		return fmt.Errorf("Failed to register custom validators, %w", err)
// 	}
// 	// return RegisterCustomValidators(validate)
// 	return nil
// }
//
// func RegisterCustomValidators(v *validator.Validate) error {
// 	err := v.RegisterValidation("not_boolean", validateNotBoolean)
// 	if err != nil {
// 		return fmt.Errorf("failed to register not_boolean validator: %w", err)
// 	}
// 	return nil
// }

// func (f *ISOStandardForm) Validate() *custom_errors.CustomError {
// 	// Empty string
// 	if f.Name == "" {
// 		return custom_errors.EmptyField("string", "name")
// 	}
//
// 	if err := validate.Struct(f); err != nil {
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
// 					return custom_errors.NewCustomError(201, "NOT Boolean error", nil)
// 				}
// 			}
// 		}
// 		return custom_errors.NewCustomError(http.StatusInternalServerError, "Unexpected validation error", nil)
// 	}
// 	return nil
// }
