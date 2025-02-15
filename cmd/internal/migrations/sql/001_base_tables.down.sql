-- Disable foreign key checks and set proper character encoding
SET FOREIGN_KEY_CHECKS = 0;
SET NAMES utf8mb4;
SET sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- Drop materialized views
DROP TABLE IF EXISTS materialized_queries;

-- Drop linking/junction tables with their indexes and constraints
DROP TABLE IF EXISTS audit_questions_comments;
DROP TABLE IF EXISTS audit_question_findings;
DROP TABLE IF EXISTS finding_corrective_actions;
DROP TABLE IF EXISTS audit_support_auditors;
DROP TABLE IF EXISTS evidence_provided;
DROP TABLE IF EXISTS audit_plan_requirements;

-- Drop core business tables with their indexes and constraints
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS findings;
DROP TABLE IF EXISTS evidence;
DROP TABLE IF EXISTS audit_questions;
DROP TABLE IF EXISTS audits;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS requirement;
DROP TABLE IF EXISTS requirement_level;
DROP TABLE IF EXISTS audit_plans;

-- Drop reference tables with their indexes and constraints
DROP TABLE IF EXISTS standards;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS reference_values;
DROP TABLE IF EXISTS reference_types;

-- Re-enable foreign key checks
SET FOREIGN_KEY_CHECKS = 1;
