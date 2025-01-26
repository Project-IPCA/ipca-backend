SET @dbname = DATABASE();
SET @tablename = "departments";
SET @columnname = "name_en";
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (table_name = @tablename)
      AND (table_schema = @dbname)
      AND (column_name = @columnname)
  ) < 1,
  "SELECT 1",
  CONCAT("ALTER TABLE ", @tablename, " DROP ", @columnname)
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;