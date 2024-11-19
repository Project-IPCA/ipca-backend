package constants

type ChannelStruct struct {
	OnlineStudent string
	LoginRepeat   string
	Log           string
}

var RedisChannel = ChannelStruct{
	OnlineStudent: "online-students",
	LoginRepeat:   "login-repeat",
	Log:           "logs",
}
