package constants

type ChannelStruct struct {
	OnlineStudent string
	LoginRepeat   string
}

var RedisChannel = ChannelStruct{
	OnlineStudent: "online-students",
	LoginRepeat:   "login-repeat",
}
