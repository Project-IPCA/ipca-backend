package constants

type ActionLogStruct struct {
	Login            string
	LoginRepeat      string
	Logout           string
	LogoutAllStudent string
}

var LogAction = ActionLogStruct{
	Login:            "Login",
	LoginRepeat:      "Login Repeat",
	Logout:           "Logout",
	LogoutAllStudent: "Logout All",
}

type PageLogStruct struct {
	Login          string
	ManageStudent  string
	ExerciseSubmit string
	LabExercise    string
}

var LogPage = PageLogStruct{
	Login:          "LOGIN",
	ManageStudent:  "MANAGE_STUDENT",
	ExerciseSubmit: "EXERCISE_SUBMIT",
	LabExercise:    "LAB_EXERCISE",
}
