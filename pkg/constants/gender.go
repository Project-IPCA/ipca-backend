package constants

type GenderStruct struct {
	Male   string
	Female string
	Other  string
}

var Gender = GenderStruct{
	Male:   "MALE",
	Female: "FEMALE",
	Other:  "OTHER",
}
