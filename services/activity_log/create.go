package activitylog

import (
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/utils"
	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (activityLogService *Service) Create(
	groupId *uuid.UUID,
	username string,
	remoteIp string,
	remotePort *int,
	agent *string,
	pageName string,
	action string,
) (models.ActivityLog, error) {
	logId := utils.NewULID()
	activityLog := builders.NewActivityLogBuilder().
		SetLogID(logId).
		SetTimestamp(time.Now()).
		SetGroupID(groupId).
		SetUsername(username).
		SetRemoteIP(remoteIp).
		SetRemotePort(remotePort).
		SetAgent(agent).
		SetPageName(pageName).
		SetAction(action).
		Build()
	activityLogService.DB.Create(&activityLog)
	return activityLog,nil
}
