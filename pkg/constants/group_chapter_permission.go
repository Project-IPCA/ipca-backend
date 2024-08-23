package constants

type GroupChapterPermStatusStruct struct {
	Na    string
	Ready string
	Open  string
	Close string
	Stop  string
}

var GroupChapterPermStatus = GroupChapterPermStatusStruct{
	Na:    "NA",
	Ready: "READY",
	Open:  "OPEN",
	Close: "CLOSE",
	Stop:  "STOP",
}
