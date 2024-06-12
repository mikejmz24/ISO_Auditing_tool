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
	ISOStandard    ISOStandard     `json:"iso_standard"`
	Name           string          `json:"name"`
	Team           string          `json:"team"`
	LeadAuditor    User            `json:"user"`
	AuditQuestions []AuditQuestion `json:"audit_questions"`
}

type AuditQuestion struct {
	ID               int              `json:"id"`
	AuditID          int              `json:"audit_id"`
	EvidenceProvided EvidenceProvided `json:"evidence_provided"`
	Question         Question         `json:"question"`
	Comments         []Comment        `json:"comments"`
}

type ISOStandard struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Clauses []Clause `json:"clauses"`
}

type Clause struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Sections []Section `json:"sections"`
}

type Section struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Evidence struct {
	ID       int    `json:"id"`
	Expected string `json:"expected"`
}

type EvidenceProvided struct {
	ID       int      `json:"id"`
	Evidence Evidence `json:"evidence"`
	Provided []string `json:"provided"`
}

type Comment struct {
	ID   int    `json:"id"`
	User User   `json:"user"`
	Text string `json:"text"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
