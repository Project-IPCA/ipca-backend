package responses

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupDetail struct {
	GroupID    uuid.UUID `json:"group_id"`
	Name       string    `json:"name"`
	Number     *int      `json:"number"`
	Day        string    `json:"day"`
	TimeStart  string    `json:"time_start"`
	TimeEnd    string    `json:"time_end"`
	Year       int       `json:"year"`
	Semester   int       `json:"semester"`
	Instructor string    `json:"instructor"`
}

type UserInfo struct {
	UserID      uuid.UUID    `json:"user_id"`
	Avatar      *string      `json:"avatar"`
	FirstName   *string      `json:"f_name"`
	LastName    *string      `json:"l_name"`
	Nickname    *string      `json:"nickname"`
	Gender      *string      `json:"gender"`
	DOB         *time.Time   `json:"dob"`
	Department  *Department  `json:"dept"`
	Email       *string      `json:"email"`
	Tel         *string      `json:"tel"`
	KmitlID     *string      `json:"kmitl_id"`
	GroupDetail *GroupDetail `json:"group_info"`
}

type Selecter struct {
	Departments []Department `json:"departments"`
}

type UserInfoResponse struct {
	UserInfo UserInfo `json:"user_info"`
	Selecter Selecter `json:"selected"`
}

func NewUserInfoResponse(user models.User, allDepartments []models.Department, classSchedule *models.ClassSchedule) *UserInfoResponse {
	deptsRes := make([]Department, 0)
	for _, dept := range allDepartments {
		deptsRes = append(deptsRes, Department{
			DeptID: dept.DeptID,
			NameTH: dept.NameTH,
			NameEN: dept.NameEN,
		})
	}

	var userDept *Department
	if user.DeptID != nil && user.Dept != nil {
		userDept = &Department{
			DeptID: *user.DeptID,
			NameTH: user.Dept.NameTH,
			NameEN: user.Dept.NameEN,
		}
	}

	var kmitlID *string
	var groupDetail *GroupDetail
	if classSchedule != nil && user.Student != nil {
		kmitlID = &user.Student.KmitlID
		groupDetail = &GroupDetail{
			GroupID:    classSchedule.GroupID,
			Name:       classSchedule.Name,
			Number:     classSchedule.Number,
			Day:        *classSchedule.Day,
			TimeStart:  *classSchedule.TimeStart,
			TimeEnd:    *classSchedule.TimeEnd,
			Semester:   *classSchedule.Semester,
			Year:       *classSchedule.Year,
			Instructor: fmt.Sprintf("%s %s", *classSchedule.Supervisor.User.FirstName, *classSchedule.Supervisor.User.LastName),
		}
	}

	return &UserInfoResponse{
		UserInfo: UserInfo{
			UserID:      user.UserID,
			Avatar:      user.Avatar,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Nickname:    user.Nickname,
			Gender:      user.Gender,
			DOB:         user.DOB,
			Email:       user.Email,
			Tel:         user.Tel,
			Department:  userDept,
			KmitlID:     kmitlID,
			GroupDetail: groupDetail,
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
			NameTH: user.Dept.NameTH,
			NameEN: user.Dept.NameEN,
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
