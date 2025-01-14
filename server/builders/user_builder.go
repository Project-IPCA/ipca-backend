package builders

import (
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type UserBuilder struct {
	UserID    uuid.UUID
	Username  string
	Password  string
	FirstName *string
	LastName  *string
	Nickname  *string
	Gender    *string
	DOB       *time.Time
	Avatar    *string
	Role      *string
	Email     *string
	Tel       *string
	Added     time.Time
	LastLogin *time.Time
	LastSeen  time.Time
	IsOnline  bool
	IsActive  bool
	AddedBy   *string
	CISession *int
	SessionID *string
	DeptID    *uuid.UUID
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (userBuilder *UserBuilder) SetID(userId uuid.UUID) (u *UserBuilder) {
	userBuilder.UserID = userId
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

func (userBuilder *UserBuilder) SetFirstName(firstName string) (u *UserBuilder) {
	userBuilder.FirstName = &firstName
	return userBuilder
}

func (userBuilder *UserBuilder) SetLastName(lastName string) (u *UserBuilder) {
	userBuilder.LastName = &lastName
	return userBuilder
}

func (userBuilder *UserBuilder) SetNickname(nickname string) (u *UserBuilder) {
	userBuilder.Nickname = &nickname
	return userBuilder
}

func (userBuilder *UserBuilder) SetGender(gender string) (u *UserBuilder) {
	userBuilder.Gender = &gender
	return userBuilder
}

func (userBuilder *UserBuilder) SetDOB(dob time.Time) (u *UserBuilder) {
	userBuilder.DOB = &dob
	return userBuilder
}

func (userBuilder *UserBuilder) SetAvatar(avatar string) (u *UserBuilder) {
	userBuilder.Avatar = &avatar
	return userBuilder
}

func (userBuilder *UserBuilder) SetRole(role string) (u *UserBuilder) {
	userBuilder.Role = &role
	return userBuilder
}

func (userBuilder *UserBuilder) SetEmail(email string) (u *UserBuilder) {
	userBuilder.Email = &email
	return userBuilder
}

func (userBuilder *UserBuilder) SetTel(tel string) (u *UserBuilder) {
	userBuilder.Tel = &tel
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
	userBuilder.CISession = &ciSession
	return userBuilder
}

func (userBuilder *UserBuilder) SetSessionID(sessionId string) (u *UserBuilder) {
	userBuilder.SessionID = &sessionId
	return userBuilder
}

func (userBuilder *UserBuilder) SetDeptID(deptId uuid.UUID) (u *UserBuilder) {
	userBuilder.DeptID = &deptId
	return userBuilder
}

func (userBuilder *UserBuilder) Build() models.User {
	user := models.User{
		UserID:    userBuilder.UserID,
		Username:  userBuilder.Username,
		Password:  userBuilder.Password,
		FirstName: userBuilder.FirstName,
		LastName:  userBuilder.LastName,
		Nickname:  userBuilder.Nickname,
		Gender:    userBuilder.Gender,
		DOB:       userBuilder.DOB,
		Avatar:    userBuilder.Avatar,
		Role:      userBuilder.Role,
		Email:     userBuilder.Email,
		Tel:       userBuilder.Tel,
		Added:     userBuilder.Added,
		LastLogin: userBuilder.LastLogin,
		LastSeen:  userBuilder.LastSeen,
		IsOnline:  userBuilder.IsOnline,
		IsActive:  userBuilder.IsActive,
		AddedBy:   userBuilder.AddedBy,
		CISession: userBuilder.CISession,
		SessionID: userBuilder.SessionID,
		DeptID:    userBuilder.DeptID,
	}
	return user
}
