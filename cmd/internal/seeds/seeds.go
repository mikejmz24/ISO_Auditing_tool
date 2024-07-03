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
	// Add other seed functions here
}

func Seed(db *sql.DB) error {
	log.Println("Starting the seeding process...")
	for _, seedFunc := range seedFuncs {
		log.Printf("Running seed function: %T\n", seedFunc)
		if err := seedFunc(db); err != nil {
			return fmt.Errorf("failed to execute seed: %w", err)
		}
	}
	log.Println("Database tables seeded successfully")
	return nil
}

func seedISOStandards(db *sql.DB) error {
	log.Println("Seeding ISO Standards...")
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM iso_standard")
	if err := row.Scan(&count); err != nil {
		log.Printf("Error querying iso_standard count: %v\n", err)
		return err
	}
	log.Printf("ISO Standards count: %d\n", count)
	if count > 0 {
		log.Println("ISO Standards already seeded")
		return nil // Data already exists
	}

	query := `INSERT INTO iso_standard (name) VALUES 
	('ISO 9001:2015'),
	('ISO 27001:2013');`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error inserting ISO Standards: %v\n", err)
		return err
	}
	log.Println("ISO Standards seeded successfully")
	return nil
}

func seedClauses(db *sql.DB) error {
	log.Println("Seeding Clauses...")
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM clause")
	if err := row.Scan(&count); err != nil {
		log.Printf("Error querying clause count: %v\n", err)
		return err
	}
	log.Printf("Clauses count: %d\n", count)
	if count > 0 {
		log.Println("Clauses already seeded")
		return nil // Data already exists
	}

	query := `INSERT INTO clause (id, iso_standard_id, name) VALUES
	(1, 1, '4. Context of the Organization'),
	(2, 1, '5 Leadership'),
	(3, 1, '6 Planning'),
	(4, 1, '7 Support'),
	(5, 1, '8 Operation'),
	(6, 1, '9 Performance evaluation'),
	(7, 1, '10 Improvement');`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error inserting Clauses: %v\n", err)
		return err
	}
	log.Println("Clauses seeded successfully")
	return nil
}

func seedSections(db *sql.DB) error {
	log.Println("Seeding Sections...")
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM section")
	if err := row.Scan(&count); err != nil {
		log.Printf("Error querying section count: %v\n", err)
		return err
	}
	log.Printf("Sections count: %d\n", count)
	if count > 0 {
		log.Println("Sections already seeded")
		return nil // Data already exists
	}

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
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error inserting Sections: %v\n", err)
		return err
	}
	log.Println("Sections seeded successfully")
	return nil
}

func seedSubsections(db *sql.DB) error {
	log.Println("Seeding Subsections...")
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM subsection")
	if err := row.Scan(&count); err != nil {
		log.Printf("Error querying subsection count: %v\n", err)
		return err
	}
	log.Printf("Subsections count: %d\n", count)
	if count > 0 {
		log.Println("Subsections already seeded")
		return nil // Data already exists
	}

	query := `INSERT INTO subsection (id, section_id, name) VALUES
	(1, 5, '5.1.1 General'),
	(2, 5, '5.1.2 Customer focus'),
	(3, 6, '5.2.1 Establishing the quality policy'),
	(4, 6, '5.2.2 Communicating the quality policy'),
	(5, 11, '7.1.1 General'),
	(5, 11, '7.1.2 People'),
	(6, 11, '7.1.3 Infrastructure'),
	(7, 11, '7.1.4 Environment for the operation of processes'),
	(8, 11, '7.1.5 Monitoring and measuring resources'),
	(9, 11, '7.1.5.1 General'),
	(10, 11, '7.1.5.2 Measurement traceability'),
	(11, 11, '7.1.6 Organizational knowledge'),
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
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error inserting Subsections: %v\n", err)
		return err
	}
	log.Println("Subsections seeded successfully")
	return nil
}
