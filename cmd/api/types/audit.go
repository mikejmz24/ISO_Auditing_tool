// This file was generated from JSON Schema using quicktype, do not modify it directly. To parse and unparse this JSON data, add this code to your project and do:
//
//    audit, err := UnmarshalAudit(bytes)
//    bytes, err = audit.Marshal()

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
	ID          int             `json:"id"`
	Datetime    time.Time       `json:"datetime"`
	Name        string          `json:"name"`
	Team        string          `json:"team"`
	LeadAuditor User            `json:"user"`
	Audit       []AuditQuestion `json:"audit"`
}

type AuditQuestion struct {
	Clause   Clause    `json:"clause"`
	Evidence Evidence  `json:"evidence"`
	Comments []Comment `json:"comments"`
}

type Clause struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Section string `json:"section"`
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

type Evidence struct {
	ID       int    `json:"id"`
	Expected string `json:"expected"`
	Provided string `json:"provided"`
}
