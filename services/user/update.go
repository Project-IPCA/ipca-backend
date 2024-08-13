package user

import (
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
)

func (userService *Service) UpdateLoginSuccess(user *models.User) error {
	timeNow := time.Now()
	user.LastLogin = &timeNow
	user.LastSeen = timeNow
	user.IsOnline = true
	userService.DB.Save(user)
	return nil
}

func (userService *Service) UpdateIsOnline(user *models.User, isOnline bool) error {
	user.IsOnline = isOnline
	userService.DB.Save(user)
	return nil
}
