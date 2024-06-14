package types

import (
	"encoding/json"
	"time"
)

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
	ID                 int              `json:"id"`
	AuditID            int              `json:"audit_id"`
	EvidenceProvidedID int              `json:"evidence_provided_id"`
	QuestionID         int              `json:"question_id"`
	EvidenceProvided   EvidenceProvided `json:"evidence_provided"`
	Question           Question         `json:"question"`
	Comments           []Comment        `json:"comments"`
}

type ISOStandard struct {
	ID      int      `form:"id" json:"id"`
	Name    string   `form:"name" json:"name"`
	Clauses []Clause `form:"clauses" json:"clauses"`
}

type Clause struct {
	ID            int       `json:"id"`
	ISOStandardID int       `json:"iso_standard_id"`
	Name          string    `json:"name"`
	Sections      []Section `json:"sections"`
}

type Section struct {
	ID        int        `json:"id"`
	ClauseID  int        `json:"clause_id"`
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID        int    `json:"id"`
	SectionID int    `json:"section_id"`
	Text      string `json:"text"`
}

type Evidence struct {
	ID       int    `json:"id"`
	Expected string `json:"expected"`
}

type EvidenceProvided struct {
	ID         int      `json:"id"`
	EvidenceID int      `json:"evidence_id"`
	Provided   []string `json:"provided"`
	Evidence   Evidence `json:"evidence"`
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
