package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

var queries = []string{
	`CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS iso_standard (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS clause (
			id INT AUTO_INCREMENT PRIMARY KEY,
			iso_standard_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (iso_standard_id) REFERENCES iso_standard(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS section (
			id INT AUTO_INCREMENT PRIMARY KEY,
			clause_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (clause_id) REFERENCES clause(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS subsection (
			id INT AUTO_INCREMENT PRIMARY KEY,
			section_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (section_id) REFERENCES section(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS question (
			id INT AUTO_INCREMENT PRIMARY KEY,
			section_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			FOREIGN KEY (section_id) REFERENCES section(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS evidence (
			id INT AUTO_INCREMENT PRIMARY KEY,
			expected VARCHAR(255) NOT NULL
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS evidence_provided (
			id INT AUTO_INCREMENT PRIMARY KEY,
			evidence_id INT NOT NULL,
			provided TEXT NOT NULL,
			FOREIGN KEY (evidence_id) REFERENCES evidence(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS comment (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			text TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES user(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS audit (
			id INT AUTO_INCREMENT PRIMARY KEY,
			datetime DATETIME NOT NULL,
			iso_standard_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			team VARCHAR(255) NOT NULL,
			user_id INT NOT NULL,
			FOREIGN KEY (iso_standard_id) REFERENCES iso_standard(id),
			FOREIGN KEY (user_id) REFERENCES user(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS audit_questions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			audit_id INT NOT NULL,
			evidence_provided_id INT NOT NULL,
			question_id INT NOT NULL,
			FOREIGN KEY (audit_id) REFERENCES audit(id),
			FOREIGN KEY (evidence_provided_id) REFERENCES evidence_provided(id),
			FOREIGN KEY (question_id) REFERENCES question(id)
		) AUTO_INCREMENT=1;`,
	`CREATE TABLE IF NOT EXISTS audit_question_comments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			audit_question_id INT NOT NULL,
			comment_id INT NOT NULL,
			FOREIGN KEY (audit_question_id) REFERENCES audit_questions(id),
			FOREIGN KEY (comment_id) REFERENCES comment(id)
		) AUTO_INCREMENT=1;`,
}

func Migrate(db *sql.DB) error {
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}
	log.Println("Database tables created successfully")
	return nil
}