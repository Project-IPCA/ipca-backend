ALTER TABLE `lab_exercises`
DROP FOREIGN KEY `fk_lab_class_infos_supervisors`,
ADD CONSTRAINT `fk_lab_class_infos_users`
FOREIGN KEY (`created_by`) REFERENCES `users` (`user_id`) 
ON DELETE CASCADE 
ON UPDATE CASCADE;

