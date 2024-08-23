package constants

type GroupAssignmentChapterItemStatusStruct struct {
	Ready  string
	Closed string
	Stop   string
	Open   string
}

var GroupAssignmentChapterItemStatus = GroupAssignmentChapterItemStatusStruct{
	Ready:  "READY",
	Closed: "CLOSED",
	Stop:   "STOP",
	Open:   "OPEN",
}
