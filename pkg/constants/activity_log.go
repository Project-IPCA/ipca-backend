package constants

type ActionLogStruct struct {
	Login       string
	LoginRepeat string
	Logout      string
}

var LogAction = ActionLogStruct{
	Login:       "LOGIN",
	LoginRepeat: "LOGIN_REPEAT",
	Logout:      "LOGOUT",
}

type PageLogStruct struct {
	Login string
}

var LogPage = PageLogStruct{
	Login: "LOGIN",
}
