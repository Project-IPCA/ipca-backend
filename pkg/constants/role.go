package constants

type RoleStruct struct {
	Student    string
	Supervisor string
	Ta         string
	Executive  string
	Beyonder   string
}

var Role = RoleStruct{
	Student:    "STUDENT",
	Supervisor: "SUPERVISOR",
	Ta:         "TA",
	Executive:  "EXECUTIVE",
	Beyonder:   "BEYONDER",
}
