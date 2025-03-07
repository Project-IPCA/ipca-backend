ALTER TABLE `lab_exercises`
DROP FOREIGN KEY `fk_lab_class_infos_users`,
ADD CONSTRAINT `fk_lab_class_infos_supervisors`
FOREIGN KEY (`created_by`) REFERENCES `supervisors` (`supervisor_id`) 
ON DELETE CASCADE 
ON UPDATE CASCADE;
