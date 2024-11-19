package constants

type ActionLogStruct struct {
	Login            string
	LoginRepeat      string
	Logout           string
	LogoutAllStudent string
}

var LogAction = ActionLogStruct{
	Login:            "LOGIN",
	LoginRepeat:      "LOGIN_REPEAT",
	Logout:           "LOGOUT",
	LogoutAllStudent: "LOGOUT_ALL",
}

type PageLogStruct struct {
	Login          string
	ManageStudent  string
	ExerciseSubmit string
}

var LogPage = PageLogStruct{
	Login:          "LOGIN",
	ManageStudent:  "MANAGE_STUDENT",
	ExerciseSubmit: "EXERCISE_SUBMIT",
}
