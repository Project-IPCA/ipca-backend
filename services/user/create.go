package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (userService *Service) Create(stuId string, role string) (userID uuid.UUID, err error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(stuId),
		bcrypt.DefaultCost,
	)
	userID = uuid.New()
	if err != nil {
		return uuid.Nil, err
	}

	user := builders.NewUserBuilder().
		SetID(userID).
		SetUsername(stuId).
		SetPassword(string(encryptedPassword)).
		SetRole(role).
		Build()

	userService.DB.Create(&user)
	return userID, nil
}
