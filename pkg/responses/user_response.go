package responses

import (
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type UserInfo struct {
	Avatar     *string     `json:"avatar"`
	FirstName  *string     `json:"f_name"`
	LastName   *string     `json:"l_name"`
	Nickname   *string     `json:"nickname"`
	Gender     *string     `json:"gender"`
	DOB        *time.Time  `json:"dob"`
	Department *Department `json:"dept"`
	Email      *string     `json:"email"`
	Tel        *string     `json:"tel"`
}

type Selecter struct {
	Departments []Department `json:"departments"`
}

type UserInfoResponse struct {
	UserInfo UserInfo `json:"user_info"`
	Selecter Selecter `json:"selected"`
}

func NewUserInfoResponse(user models.User, allDepartments []models.Department) *UserInfoResponse {
	deptsRes := make([]Department, 0)
	for _, dept := range allDepartments {
		deptsRes = append(deptsRes, Department{
			DeptID: dept.DeptID,
			Name:   dept.Name,
		})
	}

	var userDept *Department
	if user.DeptID != nil && user.Dept != nil {
		userDept = &Department{
			DeptID: *user.DeptID,
			Name:   user.Dept.Name,
		}
	}

	return &UserInfoResponse{
		UserInfo: UserInfo{
			Avatar:     user.Avatar,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Nickname:   user.Nickname,
			Gender:     user.Gender,
			DOB:        user.DOB,
			Email:      user.Email,
			Tel:        user.Tel,
			Department: userDept,
		},
		Selecter: Selecter{
			Departments: deptsRes,
		},
	}
}

type GroupInfo struct {
	GroupID uuid.UUID `json:"group_id"`
	Name    string    `json:"name"`
	Number  *int      `json:"number"`
}

type UserStudentInfoResponse struct {
	UserID     uuid.UUID   `json:"user_id"`
	KmitlID    string      `json:"kmitl_id"`
	Avatar     *string     `json:"avatar"`
	FirstName  *string     `json:"f_name"`
	LastName   *string     `json:"l_name"`
	Nickname   *string     `json:"nickname"`
	Gender     *string     `json:"gender"`
	DOB        *time.Time  `json:"dob"`
	Department *Department `json:"dept"`
	Email      *string     `json:"email"`
	Tel        *string     `json:"tel"`
	CanSubmit  bool        `json:"can_submit"`
	IsOnline   bool        `json:"is_online"`
	Group      *GroupInfo  `json:"group"`
}

func NewUserStudentInfoResponse(user models.User) *UserStudentInfoResponse {
	var userDept *Department
	if user.DeptID != nil && user.Dept != nil {
		userDept = &Department{
			DeptID: *user.DeptID,
			Name:   user.Dept.Name,
		}
	}

	var userGroup *GroupInfo
	if user.Student.GroupID != nil {
		userGroup = &GroupInfo{
			GroupID: *user.Student.GroupID,
			Name:    user.Student.Group.Name,
			Number:  user.Student.Group.Number,
		}
	}

	return &UserStudentInfoResponse{
		UserID:     user.UserID,
		KmitlID:    user.Student.KmitlID,
		Avatar:     user.Avatar,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Nickname:   user.Nickname,
		Gender:     user.Gender,
		DOB:        user.DOB,
		Email:      user.Email,
		Tel:        user.Tel,
		Department: userDept,
		CanSubmit:  user.Student.CanSubmit,
		IsOnline:   user.IsOnline,
		Group:      userGroup,
	}
}
