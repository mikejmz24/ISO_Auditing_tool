SET FOREIGN_KEY_CHECKS = 0;

SELECT CONCAT('DROP TABLE IF EXISTS `', GROUP_CONCAT(table_name SEPARATOR '`, `'), '`;')
INTO @dropStatement
FROM information_schema.tables
WHERE
    table_schema = DATABASE()
    AND table_type = 'BASE TABLE';

PREPARE stmt FROM @dropStatement;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET FOREIGN_KEY_CHECKS = 1;
