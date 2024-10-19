package groupchapterpermission

import "github.com/Project-IPCA/ipca-backend/models"

func (groupChapterPermissionService *Service) UpdateByModel(
	groupChapterPermission *models.GroupChapterPermission,
) error {
	if err := groupChapterPermissionService.DB.Save(&groupChapterPermission).Error; err != nil {
		return err
	}
	return nil
}
