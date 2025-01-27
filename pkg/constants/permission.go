package constants

type AccessTypeStruct struct {
	Deny        string
	Always      string
	Timer       string
	TimerPaused string
	DateTime    string
}

var AccessType = AccessTypeStruct{
	Deny:        "DENY",
	Always:      "ALWAYS",
	Timer:       "TIMER",
	TimerPaused: "TIMER_PASUED",
	DateTime:    "DATETIME",
}

type PermissionTypeStruct struct {
	GroupAdmin     string
	ExerciseAdmin  string
	StudentAdmin   string
	DashboardAdmin string
}

var PermissionType = PermissionTypeStruct{
	GroupAdmin:     "GROUP_ADMIN",
	ExerciseAdmin:  "EXERCISE_ADMIN",
	StudentAdmin:   "STUDENT_ADMIN",
	DashboardAdmin: "DASHBOARD_ADMIN",
}
