package constants

type ActionTypeStruct struct{
	Always string
	Deny string
	Timer string
	TimerPaused string
	DateTime string
}

var ActionType = AccessTypeStruct {
	Always: "ALWAYS",
	Deny: "DENY",
	Timer: "TIMER",
	TimerPaused: "TIME_PAUSED",
	DateTime: "DATETIME",
}

var ActionTypeList = []string{"ALWAYS","DENY", "TIMER","TIME_PAUSED","DATETIME"}

type PermissionPrefixStruct struct{
	Submit string
	Access string
}

var PermissionPrefix = PermissionPrefixStruct{
	Submit: "SUBMIT",
	Access: "ACCESS",
}