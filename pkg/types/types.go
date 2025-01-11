package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
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
	// ID      int       `json:"id" form:"id"`
	// Name    string    `json:"name" form:"name"`
	// Clauses []*Clause `json:"clauses" form:"clauses"`
	ID      int       `form:"id"`
	Name    string    `form:"name"`
	Clauses []*Clause `form:"clauses,omitempty"`
}

type Clause struct {
	ID            int        `json:"id"`
	ISOStandardID int        `json:"iso_standard_id"`
	Name          string     `json:"name"`
	Sections      []*Section `json:"sections,omitempty"`
}

type ClauseForm struct {
	ID            int        `json:"id" form:"id"`
	ISOStandardID int        `json:"iso_standard_id" form:"iso_standard_id"`
	Name          string     `json:"name" form:"name" binding:"required"`
	Sections      []*Section `json:"sections" form:"sections"`
}

type Section struct {
	ID        int         `json:"id"`
	ClauseID  int         `json:"clause_id"`
	Name      string      `json:"name"`
	Questions []*Question `json:"questions,omitempty"`
}

type SectionForm struct {
	ID        int         `json:"id" form:"id"`
	ClauseID  int         `json:"clause_id" form:"clause_id"`
	Name      string      `json:"name" form:"name" binding:"required"`
	Questions []*Question `json:"questions" form:"questions"`
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

type Evidence struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Expected   string `json:"expected"`
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
	return &ISOStandard{
		ID:      f.ID,
		Name:    f.Name,
		Clauses: f.Clauses,
	}
}

func (f *ISOStandardForm) FromISOStandard(iso *ISOStandard) {
	f.ID = iso.ID
	f.Name = iso.Name
	f.Clauses = iso.Clauses
}
