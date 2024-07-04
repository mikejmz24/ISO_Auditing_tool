package seeds

import (
	"database/sql"
	"fmt"
	"log"
)

var seedFuncs = []func(*sql.DB) error{
	seedISOStandards,
	seedClauses,
	seedSections,
	seedSubsections,
	seedQuestions,
}

func Seed(db *sql.DB) error {
	for _, seedFunc := range seedFuncs {
		if err := seedFunc(db); err != nil {
			log.Printf("Failed to execute seed function: %v", err)
			return fmt.Errorf("failed to execute seed: %w", err)
		}
	}
	log.Println("Database tables seeded successfully")
	return nil
}

func seeding(db *sql.DB, table_name string, query string) error {
	log.Printf("Seeding %v table...", table_name)
	var count int
	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %v", table_name)
	row := db.QueryRow(queryString)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error counting rows in %v: %v", table_name, err)
		return err
	}
	log.Printf("Count in %v: %d", table_name, count)
	if count > 0 {
		log.Printf("%v table already seeded, skipping...", table_name)
		return nil // Data already exists
	}
	if _, err := db.Exec(query); err != nil {
		log.Printf("Error inserting into %v: %v", table_name, err)
		return err
	}
	log.Printf("%v table seeded successfully", table_name)
	return nil
}

func seedISOStandards(db *sql.DB) error {
	query := `INSERT INTO iso_standard (name) VALUES 
	('ISO 9001:2015'),
	('ISO 27001:2013');`
	return seeding(db, "iso_standard", query)
}

func seedClauses(db *sql.DB) error {
	query := `INSERT INTO clause (id, iso_standard_id, name) VALUES
	(1, 1, '4. Context of the Organization'),
	(2, 1, '5 Leadership'),
	(3, 1, '6 Planning'),
	(4, 1, '7 Support'),
	(5, 1, '8 Operation'),
	(6, 1, '9 Performance evaluation'),
	(7, 1, '10 Improvement');`
	return seeding(db, "clause", query)
}

func seedSections(db *sql.DB) error {
	query := `INSERT INTO section (id, clause_id, name) VALUES
  (1, 1, '4.1 Understanding the organization and its context'),
	(2, 1, '4.2 Understanding the needs and expectations of interested parties'),
	(3, 1, '4.3 Determining the scope of the quality management system'),
	(4, 1, '4.4 Quality management system and its processes'),
	(5, 2, '5.1 Leadership and commitment'),
	(6, 2, '5.2 Policy'),
	(7, 2, '5.3 Organizational roles, responsibilities and authorities'),
	(8, 3, '6.1 Actions to address risks and opportunities'),
	(9, 3, '6.2 Quality objectives and planning to achieve them'),
	(10, 3, '6.3 Planning of changes'),
	(11, 4, '7.1 Resources'),
	(12, 4, '7.2 Competence'),
	(13, 4, '7.3 Awareness'),
	(14, 4, '7.4 Communication'),
	(15, 4, '7.5 Documented information'),
	(16, 5, '8.1 Operational planning and control'),
	(17, 5, '8.2 Requirements for products and services'),
	(18, 5, '8.3 Design and development of products and services'),
	(19, 5, '8.4 Control of externally provided processes, products and services'),
	(20, 5, '8.5 Production and service provision'),
	(21, 5, '8.6 Release of products and services'),
	(22, 5, '8.7 Control of nonconforming outputs'),
	(23, 6, '9.1 Monitoring, measurement, analysis and evaluation'),
	(24, 6, '9.2 Internal audit'),
	(25, 6, '9.3 Management review'),
	(26, 7, '10.1 General'),
	(27, 7, '10.2 Nonconformity and corrective action'),
	(28, 7, '10.3 Continual improvement');`
	return seeding(db, "section", query)
}

func seedSubsections(db *sql.DB) error {
	query := `INSERT INTO subsection (id, section_id, name) VALUES
	(1, 5, '5.1.1 General'),
	(2, 5, '5.1.2 Customer focus'),
	(3, 6, '5.2.1 Establishing the quality policy'),
	(4, 6, '5.2.2 Communicating the quality policy'),
	(5, 11, '7.1.1 General'),
	(6, 11, '7.1.2 People'),
	(7, 11, '7.1.3 Infrastructure'),
	(8, 11, '7.1.4 Environment for the operation of processes'),
	(9, 11, '7.1.5 Monitoring and measuring resources'),
	(10, 11, '7.1.5.1 General'),
	(11, 11, '7.1.5.2 Measurement traceability'),
	(12, 11, '7.1.6 Organizational knowledge'),
	(13, 15, '7.5.1 General'),
	(14, 15, '7.5.2 Creating and updating'),
	(15, 15, '7.5.3 Control of documented information'),
	(16, 17, '8.2.1 Customer communication'),
	(17, 17, '8.2.2 Determining the requirements for products and services'),
	(18, 17, '8.2.3 Review of the requirements for products and services'),
	(19, 17, '8.2.4 Changes to requirements for products and services'),
	(20, 18, '8.3.1 General'),
	(21, 18, '8.3.2 Design and development planning'),
	(22, 18, '8.3.3 Design and development inputs'),
	(23, 18, '8.3.4 Design and development controls'),
	(24, 18, '8.3.5 Design and development outputs'),
	(25, 18, '8.3.6 Design and development changes'),
	(26, 19, '8.4.1 General'),
	(27, 19, '8.4.2 Type and extent of control'),
	(28, 19, '8.4.3 Information for external providers'),
	(29, 20, '8.5.1 Control of production and service provision'),
	(30, 20, '8.5.2 Identification and traceability'),
	(31, 20, '8.5.3 Property belonging to customers or external providers'),
	(32, 20, '8.5.4 Preservation'),
	(33, 20, '8.5.5 Post-delivery activities'),
	(34, 20, '8.5.6 Control of changes'),
	(35, 23, '9.1.1 General'),
	(36, 23, '9.1.2 Customer satisfaction'),
	(37, 23, '9.1.3 Analysis and evaluation'),
	(38, 25, '9.3.1 General'),
	(39, 25, '9.3.2 Management review inputs'),
	(40, 25, '9.3.3 Management review outputs');`
	return seeding(db, "subsection", query)
}

func seedQuestions(db *sql.DB) error {
	query := `INSERT INTO question (id, section_id, subsection_id, name) VALUES
  (1, 3, NULL, 'What is the scope of the QMS?'),
	(2, 7, NULL, 'What roles and responsibilities are part of your team?'),
	(3, 7, NULL, 'When the team is unable to agree on a decision, who has the authority to make the final decision?'),
	(4, 7, NULL, 'What happens when a role is absent? Do you have a proces in place to determine who assumes the role responsibilities in the mean time?'),
	(5, 9, NULL, 'Define the objectives and how you apply those actions in your planning.'),
	(6, 10, NULL, 'How do you handle changes in your planning? Do you have any supporting methodology for this?'),
	(7, NULL, 5, 'Proces of TA with the manager - recruitment process and document repository'),
  (8, NULL, 6, 'How are the competences for each job determined?'),
	(9, NULL, 6, 'Determine the necessary competence for personnel performing work affecting product quality.'),
	(10, NULL, 6, 'How are you going to make sure that the training provides the missing competences?'),
	(11, NULL, 6, 'How do you handle a situation of having a new employee with a knowledge gap? For example, the colleague was hired with a certian level of knowledge but the job requires a higher level.'),
	(12, 12, NULL, 'On-boarding plans.'),
	(13, 12, NULL, 'Off-boarding plans.'),
	(14, 12, NULL, 'On-boarding for the manager.'),
	(15, 12, NULL, 'Knowledge transfer.'),
	(16, 12, NULL, 'Skill Matrix Plan?'),
	(17, 12, NULL, 'There must be evidence of removal of acceses when an employee changed teams. Check for the last change on the team and review the access removal evidence.'),
	(18, NULL, 16, 'How is communication facilitated with customers?'),
	(19, NULL, 18, 'How do you review requirements for your product or service?'),
	(20, NULL, 19, 'How do you manage changes in requirements for your product or service?'),
	(21, 18, NULL, 'What development framework does the team follow?'),
	(22, 18, NULL, 'Is there an SDP describing the team\'s development process?'),
	(23, NULL, 22, 'What is your process for gathering requirements?'),
	(24, NULL, 22, 'Where do your requirements come from?'),
	(25, NULL, 22, 'How fo you ensusre that a requirement is ready to be taken?'),
	(26, NULL, 22, 'How are Non-Functional Requirements (NFR) handled?'),
	(27, NULL, 22, 'What is you feature architecture design process?'),
	(28, NULL, 22, 'How do you handle architectural decisions?'),
	(29, NULL, 22, 'How are project risks managed?');`
	return seeding(db, "question", query)
}
