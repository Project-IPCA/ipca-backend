package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (userService *Service) CreateQuick(
	username string,
	password string,
	firstName string,
	lastName string,
	gender string,
	role string,
) (userID uuid.UUID, err error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	userID = uuid.New()
	if err != nil {
		return uuid.Nil, err
	}

	user := builders.NewUserBuilder().
		SetID(userID).
		SetUsername(username).
		SetPassword(string(encryptedPassword)).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetGender(gender).
		SetRole(role).
		Build()

	userService.DB.Create(&user)
	return userID, nil
}

func (userService *Service) Create(
	username string,
	password string,
	firstName string,
	lastName string,
	gender string,
	role string,
	deptID uuid.UUID,
) (userID uuid.UUID, err error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	userID = uuid.New()
	if err != nil {
		return uuid.Nil, err
	}

	user := builders.NewUserBuilder().
		SetID(userID).
		SetUsername(username).
		SetPassword(string(encryptedPassword)).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetGender(gender).
		SetRole(role).
		SetDeptID(deptID).
		Build()

	userService.DB.Create(&user)
	return userID, nil
}
