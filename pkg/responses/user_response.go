package responses

import (
	"time"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
)

type UserInfoResponse struct {
	Avatar     *string    `json:"avatar"`
	FirstName  *string    `json:"f_name"`
	LastName   *string    `json:"l_name"`
	Nickname   *string    `json:"nickname"`
	Gender     *string    `json:"gender"`
	DOB        *time.Time `json:"dob"`
	Department *string    `json:"dept"`
	Email      *string    `json:"email"`
	Tel        *string    `json:"tel"`
}

func NewUserInfoResponse(user models.User) *UserInfoResponse {
	var dept *string
	if user.Role != nil {
		if *user.Role == constants.Role.Supervisor {
			if user.Supervisor != nil && user.Supervisor.Dept != nil {
				dept = user.Supervisor.Dept
			}
		} else if *user.Role == constants.Role.Student {
			if user.Student != nil && user.Student.DeptID != nil {
				dept = &user.Student.Dept.Name
			}
		}
	}
	return &UserInfoResponse{
		Avatar:     user.Avatar,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Nickname:   user.Nickname,
		Gender:     user.Gender,
		DOB:        user.DOB,
		Department: dept,
		Email:      user.Email,
		Tel:        user.Tel,
	}
}
