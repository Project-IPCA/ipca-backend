package constants

type AccessTypeStruct struct {
	Deny       string
	Always string
	Timer      string
	TimerPaused string
	DateTime string
}

var AccessType = AccessTypeStruct{
	Deny:       "DENY",
	Always: "ALWAYS",
	Timer: "TIMER",
	TimerPaused: "TIMER_PASUED",
	DateTime: "DATETIME",
}