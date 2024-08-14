package builders

import (
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ActivityLogBuilder struct {
	LogID      ulid.ULID
	Timestamp  time.Time
	GroupID    *uuid.UUID
	Username   string
	RemoteIP   string
	RemotePort *int
	Agent      *string
	PageName   string
	Action     string
	CI         *uint
}

func NewActivityLogBuilder() *ActivityLogBuilder {
	return &ActivityLogBuilder{}
}

func (activityLogBuilder *ActivityLogBuilder) SetLogID(logId ulid.ULID) (al *ActivityLogBuilder) {
	activityLogBuilder.LogID = logId
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetTimestamp(
	timestamp time.Time,
) (al *ActivityLogBuilder) {
	activityLogBuilder.Timestamp = timestamp
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetGroupID(
	groupId *uuid.UUID,
) (al *ActivityLogBuilder) {
	activityLogBuilder.GroupID = groupId
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetUsername(
	username string,
) (al *ActivityLogBuilder) {
	activityLogBuilder.Username = username
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetRemoteIP(
	remoteIp string,
) (al *ActivityLogBuilder) {
	activityLogBuilder.RemoteIP = remoteIp
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetRemotePort(
	remotePort *int,
) (al *ActivityLogBuilder) {
	activityLogBuilder.RemotePort = remotePort
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetAgent(
	agent *string,
) (al *ActivityLogBuilder) {
	activityLogBuilder.Agent = agent
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetPageName(
	pageName *string,
) (al *ActivityLogBuilder) {
	activityLogBuilder.PageName = *pageName
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetAction(
	action string,
) (al *ActivityLogBuilder) {
	activityLogBuilder.Action = action
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) SetCI(
	ci *uint,
) (al *ActivityLogBuilder) {
	activityLogBuilder.CI = ci
	return activityLogBuilder
}

func (activityLogBuilder *ActivityLogBuilder) Build() models.ActivityLog {
	activityLog := models.ActivityLog{
		LogID:      activityLogBuilder.LogID,
		Timestamp:  activityLogBuilder.Timestamp,
		GroupID:    activityLogBuilder.GroupID,
		Username:   activityLogBuilder.Username,
		RemoteIP:   activityLogBuilder.RemoteIP,
		RemotePort: activityLogBuilder.RemotePort,
		Agent:      activityLogBuilder.Agent,
		PageName:   activityLogBuilder.PageName,
		Action:     activityLogBuilder.Action,
		CI:         activityLogBuilder.CI,
	}
	return activityLog
}
