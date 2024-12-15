package constants

type ChannelStruct struct {
	OnlineStudent   string
	UserEvent       string
	LoginRepeat     string
	Log             string
	GroupPermission string
}

var RedisChannel = ChannelStruct{
	OnlineStudent:   "online-students",
	UserEvent:       "user-event",
	LoginRepeat:     "login-repeat",
	Log:             "logs",
	GroupPermission: "group-permission",
}
