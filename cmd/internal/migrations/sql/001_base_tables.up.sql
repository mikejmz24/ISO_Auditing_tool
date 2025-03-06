-- Enable strict mode and proper character encoding
SET sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- Create reference_types table
CREATE TABLE IF NOT EXISTS reference_types (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key'
    , `name` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Type name referencing table fields (e.g., finding_type, audit_status)'
    , `description` VARCHAR(255) COMMENT 'Description of the reference type'
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp'
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp'
) ENGINE = InnoDB COMMENT = 'Master table for reference data categories';

-- Create reference_values table
CREATE TABLE IF NOT EXISTS reference_values (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key'
    , type_id INT NOT NULL COMMENT 'Reference to reference_types.id'
    , `code` VARCHAR(50) NOT NULL COMMENT 'Short code for the value'
    , `name` VARCHAR(100) NOT NULL COMMENT 'Display name'
    , `description` VARCHAR(255) COMMENT 'Description of the value'
    , is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Whether this value is currently active'
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , deleted_at TIMESTAMP NULL COMMENT 'Soft delete timestamp'
    , CONSTRAINT fk_reference_values_type FOREIGN KEY (type_id) REFERENCES reference_types (id)
    , UNIQUE KEY uk_reference_value (type_id, code)
    , INDEX idx_reference_values_active (is_active)
    , INDEX idx_reference_values_type (type_id)
) ENGINE = InnoDB COMMENT = 'Stores all reference/lookup values';

-- Create audit_plan_requirements table
CREATE TABLE IF NOT EXISTS audit_plan_requirements (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_plan_id INT NOT NULL
    , requirement VARCHAR(255) NOT NULL
    , CONSTRAINT fk_audit_plan_requirements_audit_plan FOREIGN KEY (audit_plan_id) REFERENCES audit_plans (id)
    , INDEX idx_audit_plan_requirements_plan (audit_plan_id)
) ENGINE = InnoDB COMMENT = 'Stores specific requirements (checklist) for audit plans';

-- Create audit_plans table
CREATE TABLE IF NOT EXISTS audit_plans (
    id INT AUTO_INCREMENT PRIMARY KEY
    , standard_id INT NOT NULL
    , lead_auditor_id INT NOT NULL
    , `name` VARCHAR(255) NOT NULL
    , status_id INT NOT NULL COMMENT 'Reference to reference_values.id for status'
    , scheduled_date DATETIME NOT NULL
    , team VARCHAR(255) NOT NULL
    , `scope` VARCHAR(255) NOT NULL
    , type_id INT NOT NULL COMMENT 'Reference to reference_values.id for audit types'
    , is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Whether this value is currently active'
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , deleted_at TIMESTAMP NULL
    , CONSTRAINT fk_audit_plans_standard FOREIGN KEY (standard_id) REFERENCES standards (id)
    , CONSTRAINT fk_audit_plans_lead_auditor FOREIGN KEY (lead_auditor_id) REFERENCES users (id)
    , CONSTRAINT fk_audit_plans_status FOREIGN KEY (status_id) REFERENCES reference_values (id)
    , CONSTRAINT fk_audit_plans_type FOREIGN KEY (type_id) REFERENCES reference_values (id)
    , INDEX idx_audit_plans_standard (standard_id)
    , INDEX idx_audit_plans_user (lead_auditor_id)
    , INDEX idx_audit_plans_status (status_id)
    , INDEX idx_audit_plans_type (type_id)
    , INDEX idx_audit_plans_scheduled (scheduled_date)
    , INDEX idx_audit_plans_active (deleted_at)
) ENGINE = InnoDB COMMENT = 'Stores audit plans';

-- Create users table with improved structure
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key'
    , email VARCHAR(255) NOT NULL COMMENT 'User email address'
    , `name` VARCHAR(255) NOT NULL COMMENT 'User full name'
    , role_id INT NOT NULL COMMENT 'Reference to reference_values.id for user roles'
    , is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Whether user account is active'
    , last_login_at TIMESTAMP NULL COMMENT 'Last successful login timestamp'
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , deleted_at TIMESTAMP NULL COMMENT 'Soft delete timestamp'
    , CONSTRAINT uk_users_email UNIQUE (email)
    , CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES reference_values (id)
    , INDEX idx_users_role (role_id)
    , INDEX idx_users_active (is_active, deleted_at)
) ENGINE = InnoDB COMMENT = 'Stores user account information';

-- Create audit_support_auditors table
CREATE TABLE IF NOT EXISTS audit_support_auditors (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_id INT NOT NULL
    , user_id INT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , CONSTRAINT fk_audit_support_auditors_audit FOREIGN KEY (audit_id) REFERENCES audit_plans (id)
    , CONSTRAINT fk_audit_support_auditors_user FOREIGN KEY (user_id) REFERENCES users (id)
    , INDEX idx_audit_support_auditors_audit_id (audit_id)
    , INDEX idx_audit_support_auditors_user_id (user_id)
) ENGINE = InnoDB COMMENT = 'Stores support auditors for each audit';

-- Create audits table
CREATE TABLE IF NOT EXISTS audits (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_plans_id INT NOT NULL
    , CONSTRAINT fk_audits_audit_plan FOREIGN KEY (audit_plans_id) REFERENCES audit_plans (id)
    , INDEX idx_audits_audit_plans (audit_plans_id)
) ENGINE = InnoDB COMMENT = 'Links audit plans with audit_questions and audit_support_auditors';

-- Create standards table
CREATE TABLE IF NOT EXISTS standards (
    id INT AUTO_INCREMENT PRIMARY KEY
    , `name` VARCHAR(100) NOT NULL
    , `description` VARCHAR(255)
    , `version` VARCHAR(50) NOT NULL
) ENGINE = InnoDB COMMENT = 'Stores audit standards information';

-- Create requirement_level table
CREATE TABLE IF NOT EXISTS requirement_level (
    id INT AUTO_INCREMENT PRIMARY KEY
    , standard_id INT NOT NULL
    , level_name VARCHAR(255) NOT NULL COMMENT 'The naming convention for each level. e.g., Clause, Subclause, Requirement, etc.'
    , level_order INT NOT NULL COMMENT 'The level of the item. Nested location'
    , CONSTRAINT fk_requirement_level_standard FOREIGN KEY (standard_id) REFERENCES standards (id)
    , INDEX idx_requirements_level_standard (standard_id)
) ENGINE = InnoDB COMMENT = 'Stores the requirement level for standards allowing different schemas';

-- Create requirement table
CREATE TABLE IF NOT EXISTS requirement (
    id INT AUTO_INCREMENT PRIMARY KEY
    , standard_id INT NOT NULL
    , level_id INT NOT NULL COMMENT 'The level of the item. Nested location'
    , parent_id INT NULL COMMENT 'The parent id of the requirement for self joining. Can be null'
    , reference_code VARCHAR(50) NOT NULL COMMENT 'The number of the standard item. e.g., 4.1, 2.3.4, etc'
    , `name` VARCHAR(255) NOT NULL COMMENT 'The given name. e.g., Planning'
    , `description` TEXT NULL
    , CONSTRAINT fk_requirement_standard FOREIGN KEY (standard_id) REFERENCES standards (id)
    , INDEX idx_requirements_level_standard (standard_id)
) ENGINE = InnoDB COMMENT = 'Stores the requirement level for standards allowing different schemas';

-- Create questions table
CREATE TABLE IF NOT EXISTS questions (
    id INT AUTO_INCREMENT PRIMARY KEY
    , requirement_id INT NOT NULL
    , question VARCHAR(255) NOT NULL
    , guidance TEXT
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , CONSTRAINT fk_questions_requirement FOREIGN KEY (requirement_id) REFERENCES requirement (id)
    , INDEX idx_questions_requirement (requirement_id)
) ENGINE = InnoDB COMMENT = 'Stores audit questions';

-- Create evidence table
CREATE TABLE IF NOT EXISTS evidence (
    id INT AUTO_INCREMENT PRIMARY KEY
    , question_id INT NOT NULL
    , type_id INT NOT NULL COMMENT 'Reference to reference_values.id for evidence types'
    , expected TEXT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , CONSTRAINT fk_evidence_question FOREIGN KEY (question_id) REFERENCES questions (id)
    , CONSTRAINT fk_evidence_type FOREIGN KEY (type_id) REFERENCES reference_values (id)
    , INDEX idx_evidence_question (question_id)
    , INDEX idx_evidence_type (type_id)
) ENGINE = InnoDB COMMENT = 'Stores expected evidence requirements';

-- Create evidence_provided table
CREATE TABLE IF NOT EXISTS evidence_provided (
    id INT AUTO_INCREMENT PRIMARY KEY
    , evidence_id INT NOT NULL
    , user_id INT NOT NULL
    , type_id INT NOT NULL COMMENT 'Reference to reference_values.id for evidence types'
    , evidence VARCHAR(255) NOT NULL
    , retention_days INT NOT NULL
    , confidentiality_id INT NOT NULL COMMENT 'Reference to reference_values.id for confidentiality levels'
    , status_id INT NOT NULL COMMENT 'Reference to reference_values.id for status'
    , is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Whether this value is currently active'
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , deleted_at TIMESTAMP NULL
    , CONSTRAINT fk_evidence_provided_evidence FOREIGN KEY (evidence_id) REFERENCES evidence (id)
    , CONSTRAINT fk_evidence_provided_user FOREIGN KEY (user_id) REFERENCES users (id)
    , CONSTRAINT fk_evidence_provided_type FOREIGN KEY (type_id) REFERENCES reference_values (id)
    , CONSTRAINT fk_evidence_provided_confidentiality FOREIGN KEY (confidentiality_id) REFERENCES reference_values (id)
    , CONSTRAINT fk_evidence_provided_status FOREIGN KEY (status_id) REFERENCES reference_values (id)
    , INDEX idx_evidence_provided_evidence (evidence_id)
    , INDEX idx_evidence_provided_user (user_id)
    , INDEX idx_evidence_provided_type (type_id)
    , INDEX idx_evidence_provided_confidentiality (confidentiality_id)
    , INDEX idx_evidence_provided_status (status_id)
    , INDEX idx_evidence_provided_retention (retention_days)
    , INDEX idx_evidence_provided_active (deleted_at)
    , INDEX idx_evidence_composite_search (evidence_id, type_id, status_id, deleted_at)
) ENGINE = InnoDB COMMENT = 'Stores actual evidence provided';

-- Create finding_corrective_actions table
CREATE TABLE IF NOT EXISTS finding_corrective_actions (
    id INT AUTO_INCREMENT PRIMARY KEY
    , finding_id INT NOT NULL
    , `action` VARCHAR(255) NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , CONSTRAINT fk_audit_corrective_actions_finding FOREIGN KEY (finding_id) REFERENCES findings (id)
    , INDEX idx_finding_corrective_actions_finding_id (finding_id)
) ENGINE = InnoDB COMMENT = 'Stores all the corrective actions per finding';

-- Create findings table
CREATE TABLE IF NOT EXISTS findings (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_id INT NOT NULL
    , question_id INT NOT NULL
    , finding_type_id INT NOT NULL COMMENT 'Reference to reference_values.id for finding types'
    , severity_id INT NOT NULL COMMENT 'Reference to reference_values.id for severity levels'
    , `description` TEXT NOT NULL
    , due_date TIMESTAMP NOT NULL
    , responsible_user_id INT NOT NULL
    , status_id INT NOT NULL COMMENT 'Reference to reference_values.id for status'
    , is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Whether this value is currently active'
    , created_by INT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , deleted_at TIMESTAMP NULL
    , CONSTRAINT fk_findings_audit FOREIGN KEY (audit_id) REFERENCES audit_plans (id)
    , CONSTRAINT fk_findings_question FOREIGN KEY (question_id) REFERENCES questions (id)
    , CONSTRAINT fk_findings_type FOREIGN KEY (finding_type_id) REFERENCES reference_values (id)
    , CONSTRAINT fk_findings_severity FOREIGN KEY (severity_id) REFERENCES reference_values (id)
    , CONSTRAINT fk_findings_user FOREIGN KEY (responsible_user_id) REFERENCES users (id)
    , CONSTRAINT fk_findings_status FOREIGN KEY (status_id) REFERENCES reference_values (id)
    , CONSTRAINT fk_findings_creator FOREIGN KEY (created_by) REFERENCES users (id)
    , INDEX idx_findings_audit (audit_id)
    , INDEX idx_findings_question (question_id)
    , INDEX idx_findings_finding (finding_type_id)
    , INDEX idx_findings_severity (severity_id)
    , INDEX idx_findings_user (responsible_user_id)
    , INDEX idx_findings_status (status_id)
    , INDEX idx_findings_created_by (created_by)
    , INDEX idx_findings_due_date (due_date)
    , INDEX idx_findings_active (deleted_at)
    , INDEX idx_composite_search (audit_id, finding_type_id, severity_id, status_id, deleted_at)
) ENGINE = InnoDB COMMENT = 'Stores audit findings';

-- Create audit_question_findings table
CREATE TABLE IF NOT EXISTS audit_question_findings (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_question_id INT NOT NULL
    , finding_id INT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , CONSTRAINT fk_audit_question_findings_audit_question FOREIGN KEY (audit_question_id) REFERENCES audit_questions (id)
    , CONSTRAINT fk_audit_question_findings_finding FOREIGN KEY (finding_id) REFERENCES findings (id)
    , INDEX idx_audit_question_findings_audit_question (audit_question_id)
    , INDEX idx_audit_question_findings_finding (finding_id)
) ENGINE = InnoDB COMMENT = 'Links audit questions to findings';

-- Create audit_questions table
CREATE TABLE IF NOT EXISTS audit_questions (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_id INT NOT NULL
    , question_id INT NOT NULL
    , CONSTRAINT fk_audit_questions_audit FOREIGN KEY (audit_id) REFERENCES audit_plans (id)
    , CONSTRAINT fk_audit_questions_question FOREIGN KEY (question_id) REFERENCES questions (id)
    , INDEX idx_audit_questions_audit (audit_id)
    , INDEX idx_audit_questions_question (question_id)
) ENGINE = InnoDB COMMENT = 'Stores questions for specific audits';

-- Create audit_questions_comments table
CREATE TABLE IF NOT EXISTS audit_questions_comments (
    id INT AUTO_INCREMENT PRIMARY KEY
    , audit_question_id INT NOT NULL
    , comment_id INT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , CONSTRAINT fk_audit_questions_comments_audit_question FOREIGN KEY (audit_question_id) REFERENCES audit_questions (id)
    , CONSTRAINT fk_audit_questions_comments_comment FOREIGN KEY (comment_id) REFERENCES comments (id)
    , INDEX idx_audit_questions_comments_audit_question (audit_question_id)
    , INDEX idx_audit_questions_comments_comment (comment_id)
) ENGINE = InnoDB COMMENT = 'Links audit questions to comments';

-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY
    , user_id INT NOT NULL
    , `comment` TEXT NOT NULL
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users (id)
    , INDEX idx_comments_user (user_id)
) ENGINE = InnoDB COMMENT = 'Stores user comments';

-- Create materialized_queries table with improved structure
CREATE TABLE IF NOT EXISTS materialized_queries (
    id INT AUTO_INCREMENT PRIMARY KEY
    , query_name VARCHAR(100) NOT NULL
    , query_definition TEXT NOT NULL COMMENT 'Original query that generates the data'
    , `data` JSON NOT NULL
    , `version` INT NOT NULL DEFAULT 1
    , is_refreshing BOOLEAN NOT NULL DEFAULT FALSE
    , refresh_interval INT NOT NULL DEFAULT 3600 COMMENT 'Refresh interval in seconds'
    , last_refresh_at TIMESTAMP NULL
    , next_refresh_at TIMESTAMP NULL
    , refresh_started_at TIMESTAMP NULL
    , error_count INT NOT NULL DEFAULT 0
    , last_error TEXT
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , UNIQUE KEY uk_materialized_query_name (query_name)
    , INDEX idx_materialized_queries_refresh (is_refreshing, next_refresh_at)
) ENGINE = InnoDB COMMENT = 'Stores materialized query results for complex joins simplyfying interface to consume data';

CREATE TABLE IF NOT EXISTS drafts (
    id INT AUTO_INCREMENT PRIMARY KEY
    , type_id INT NOT NULL COMMENT 'Reference to reference_values.id ie standard, audit, audit_plan, finding'
    , object_id INT NULL COMMENT 'ID of existing object being edited, NULL if new object'
    , status_id INT NOT NULL COMMENT 'Reference to reference_values.id ie draft, pending_approval, rejected, published'
    , `version` INT NOT NULL DEFAULT 1 COMMENT 'Draft version number'
    , `data` JSON NOT NULL COMMENT 'Complete JSON representation of the object'
    , diff JSON NULL COMMENT 'Changes from previous version (for existing objects)'
    , user_id INT NOT NULL COMMENT 'User who created/owns this draft'
    , approver_id INT NULL COMMENT 'User who approved/rejected the draft'
    , approval_comment TEXT NULL
    , published_at TIMESTAMP NULL COMMENT 'When the draft was published'
    , publish_error TEXT NULL COMMENT 'Any error during publishing process'
    , created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    , expires_at TIMESTAMP NULL COMMENT 'Optional expiration for abandoned drafts'

    , INDEX idx_drafts_object (type_id, object_id)
    , INDEX idx_drafts_user (user_id, status)
    , INDEX idx_drafts_status (status_id, expires_at)
) ENGINE = InnoDB COMMENT = 'Stores draft JSON data before publishing to normalized tables';
