package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
)

func (userService *Service) UpdateLoginSuccess(user *models.User) {
	timeNow := time.Now()
	user.LastLogin = &timeNow
	user.LastSeen = timeNow
	user.IsOnline = true
	userService.DB.Save(user)
}

func (userService *Service) UpdateIsOnline(user *models.User, isOnline bool) {
	user.IsOnline = isOnline
	userService.DB.Save(user)
}

func (userService *Service) UpdateUserInfo(
	user *models.User,
	request *requests.UpdateUserInfoRequest,
) {
	var dob *time.Time
	if request.DOB != nil {
		dobParse, err := time.Parse("2006-01-02", *request.DOB)
		if err == nil {
			dob = &dobParse
		}
		if dob != nil {
			user.DOB = dob
		}
	}

	if request.Avatar != nil {
		user.Avatar = request.Avatar
	}

	if request.DeptID != nil {
		user.DeptID = request.DeptID
	}

	if request.Nickname != nil {
		user.Nickname = request.Nickname
	}

	if request.Gender != nil {
		user.Gender = request.Gender
	}

	if request.Email != nil {
		user.Email = request.Email
	}

	if request.Tel != nil {
		user.Tel = request.Tel
	}

	if request.NewPassword != nil {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(*request.NewPassword),
			bcrypt.DefaultCost,
		)
		if err == nil {
			user.Password = string(encryptedPassword)
		}
	}

	userService.DB.Save(user)
}

func (userService *Service) ResetUserStudentPassword(
	user *models.User,
	kmitlId string,
) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(kmitlId),
		bcrypt.DefaultCost,
	)
	if err == nil {
		user.Password = string(encryptedPassword)
		userService.DB.Save(user)
	}
}
