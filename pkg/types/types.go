package types

import (
	// "ISO_Auditing_Tool/pkg/custom_errors"
	// "encoding/json"
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
func DecodeForm(values url.Values, dst any) error {
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

type Standard struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Version      string        `json:"version"`
	Requirements []Requirement `json:"requirements"`
}

type Requirement struct {
	ID            int        `json:"id"`
	StandardID    int        `json:"standard_id"`
	LevelID       int        `json:"level_id"`
	ParentID      int        `json:"parent_id"`
	ReferenceCode string     `json:"reference_code"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Questions     []Question `json:"questions"`
}

type Question struct {
	ID            int        `json:"id"`
	RequirementID int        `json:"requirement_id"`
	Question      string     `json:"question"`
	Guidance      string     `json:"guidance"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Evidence      []Evidence `json:"evidence"`
}

type Evidence struct {
	ID         int            `json:"id"`
	QuestionID int            `json:"question_id"`
	TypeVal    ReferenceValue `json:"type"`
	Expected   string         `json:"expected"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type ReferenceValue struct {
	ID          int       `json:"id"`
	TypeID      int       `json:"type_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuditQuestion struct {
	ID               int                `json:"id"`
	AuditID          int                `json:"audit_id"`
	QuestionID       int                `json:"question_id"`
	EvidenceProvided []EvidenceProvided `json:"evidence_provided"`
	Comments         []Comment          `json:"comments"`
}

type Draft struct {
	ID              int             `json:"id"`
	TypeID          int             `json:"type_id"`
	ObjectID        int             `json:"object_id"`
	StatusID        int             `json:"status_id"`
	Version         int             `json:"version"`
	Data            json.RawMessage `json:"data"`
	Diff            json.RawMessage `json:"diff"`
	UserID          int             `json:"user_id"`
	ApproverID      int             `json:"approver_id"`
	ApprovalComment string          `json:"approval_comment"`
	PublishError    string          `json:"publish_error"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	ExpiresAt       time.Time       `json:"expires_at"`
}

type MaterializedJSONQuery struct {
	ID         int             `json:"id"`
	Name       string          `json:"query_name"`
	EntityType string          `json:"entity_type"` // standard, requirement, question, evidence, standard_full
	EntityID   int             `json:"entity_id"`
	Definition string          `json:"query_definition"` // Query definition to debug on MySQL
	Data       json.RawMessage `json:"data"`
	Version    int             `json:"version"`
	ErrorCount int             `json:"error_count"`
	LastError  string          `json:"last_error"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at"`
}

type MaterializedHTMLQuery struct {
	ID          int        `json:"id"`
	Name        string     `json:"query_name"`
	ViewPath    string     `json:"view_path"`
	HTMLContent string     `json:"html_content"`
	Version     int        `json:"version"`
	ErrorCount  int        `json:"error_count"`
	LastError   string     `json:"last_error"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ISOStandardForm struct {
	// Name    string        `form:"name" validate:"required,min=3,max=100,not_boolean"`
	Name string `form:"name" validate:"required,min=3,max=100,not_boolean"`
	// Clauses []*ClauseForm `form:"clauses,omitempty"`
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

// AuditContentModification represents a modification to audit content
type AuditContentModification struct {
	ContentType     string          `json:"content_type"`     // "requirement", "question", "evidence"
	ContentID       int             `json:"content_id"`       // ID of the content being modified
	Action          string          `json:"action"`           // "modify", "add", "remove"
	ModifiedContent json.RawMessage `json:"modified_content"` // New/modified content as JSON
	OriginalContent json.RawMessage `json:"original_content"` // Original content for reference
	ModifiedBy      int             `json:"modified_by"`      // User ID who made the change
	ModifiedAt      time.Time       `json:"modified_at"`      // When the change was made
	Reason          string          `json:"reason"`           // Optional reason for the change
}

// RequirementModification represents a modified requirement
type RequirementModification struct {
	ID                  int    `json:"id"`
	StandardID          int    `json:"standard_id"`
	LevelID             int    `json:"level_id"`
	ParentID            int    `json:"parent_id"`
	ReferenceCode       string `json:"reference_code"`
	Name                string `json:"name"`
	Description         string `json:"description"`          // This is what gets modified
	OriginalDescription string `json:"original_description"` // Keep original for reference
}

// AuditContentDraft represents the draft structure for content modifications
type AuditContentDraft struct {
	TypeID   int                      `json:"type_id"`   // Reference to "audit_content" type
	ObjectID int                      `json:"object_id"` // ID of content being modified
	StatusID int                      `json:"status_id"` // Draft status
	Version  int                      `json:"version"`
	Data     AuditContentModification `json:"data"`    // The modification details
	UserID   int                      `json:"user_id"` // Admin making the change
}

// RequirementEditForm represents the form for editing requirements
type RequirementEditForm struct {
	RequirementID int    `form:"requirement_id" validate:"required"`
	Name          string `form:"name" validate:"required,min=3,max=255"`
	Description   string `form:"description" validate:"required,min=10"`
	Reason        string `form:"reason" validate:"max=500"`
}

// Validate validates the requirement edit form
func (f *RequirementEditForm) Validate() error {
	if f.RequirementID <= 0 {
		return ErrRequired
	}
	if f.Name == "" {
		return ErrRequired
	}
	if f.Description == "" {
		return ErrRequired
	}
	return nil
}
