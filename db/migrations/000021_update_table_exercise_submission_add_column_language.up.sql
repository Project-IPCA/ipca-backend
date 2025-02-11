ALTER TABLE `exercise_submissions` 
ADD COLUMN `language` ENUM('C', 'PYTHON') NOT NULL DEFAULT 'PYTHON';
