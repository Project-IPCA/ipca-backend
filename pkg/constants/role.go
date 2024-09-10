package constants

type RoleStruct struct {
	Admin      string
	Editor     string
	Author     string
	Student    string
	Supervisor string
	Staff      string
	Ta         string
}

var Role = RoleStruct{
	Admin:      "ADMIN",
	Editor:     "EDITOR",
	Author:     "AUTHOR",
	Student:    "STUDENT",
	Supervisor: "SUPERVISOR",
	Staff:      "STAFF",
	Ta:         "TA",
}
