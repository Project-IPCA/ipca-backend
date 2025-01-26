package models

import "github.com/google/uuid"

type RolePermission struct {
	PermissionID uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Role         string    `gorm:"type:enum('STUDENT','SUPERVISOR','EXECUTIVE','TA','BEYONDER')"`
	Permission   string    `gorm:"type:enum('GROUP_ADMIN', 'EXERCISE_ADMIN', 'STUDENT_ADMIN', 'DASHBOARD_ADMIN')"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
