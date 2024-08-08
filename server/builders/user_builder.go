package builders

import (
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type UserBuilder struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Role      string
	Added     time.Time
	LastLogin *time.Time
	LastSeen  time.Time
	IsOnline  bool
	IsActive  bool
	AddedBy   *string
	CiSession *int
	SessionID *string
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (userBuilder *UserBuilder) SetID(id uuid.UUID) (u *UserBuilder) {
	userBuilder.ID = id
	return userBuilder
}

func (userBuilder *UserBuilder) SetUsername(username string) (u *UserBuilder) {
	userBuilder.Username = username
	return userBuilder
}

func (userBuilder *UserBuilder) SetPassword(password string) (u *UserBuilder) {
	userBuilder.Password = password
	return userBuilder
}

func (userBuilder *UserBuilder) SetRole(role string) (u *UserBuilder) {
	userBuilder.Role = role
	return userBuilder
}

func (userBuilder *UserBuilder) SetAdded(added time.Time) (u *UserBuilder) {
	userBuilder.Added = added
	return userBuilder
}

func (userBuilder *UserBuilder) SetLastLogin(lastLogin time.Time) (u *UserBuilder) {
	userBuilder.LastLogin = &lastLogin
	return userBuilder
}

func (userBuilder *UserBuilder) SetLastSeen(lastSeen time.Time) (u *UserBuilder) {
	userBuilder.LastSeen = lastSeen
	return userBuilder
}

func (userBuilder *UserBuilder) SetIsOnline(isOnline bool) (u *UserBuilder) {
	userBuilder.IsOnline = isOnline
	return userBuilder
}

func (userBuilder *UserBuilder) SetIsActive(isActive bool) (u *UserBuilder) {
	userBuilder.IsActive = isActive
	return userBuilder
}

func (userBuilder *UserBuilder) SetAddedBy(addedBy string) (u *UserBuilder) {
	userBuilder.AddedBy = &addedBy
	return userBuilder
}

func (userBuilder *UserBuilder) SetCiSession(ciSession int) (u *UserBuilder) {
	userBuilder.CiSession = &ciSession
	return userBuilder
}

func (userBuilder *UserBuilder) SetSessionID(sessionId string) (u *UserBuilder) {
	userBuilder.SessionID = &sessionId
	return userBuilder
}

func (userBuilder *UserBuilder) Build() models.User {
	user := models.User{
		ID:        userBuilder.ID,
		Username:  userBuilder.Username,
		Password:  userBuilder.Password,
		Role:      userBuilder.Role,
		Added:     userBuilder.Added,
		LastLogin: userBuilder.LastLogin,
		LastSeen:  userBuilder.LastSeen,
		IsOnline:  userBuilder.IsOnline,
		IsActive:  userBuilder.IsActive,
		AddedBy:   userBuilder.AddedBy,
		CiSession: userBuilder.CiSession,
		SessionID: userBuilder.SessionID,
	}
	return user
}
