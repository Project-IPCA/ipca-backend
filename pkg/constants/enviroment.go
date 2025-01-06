package constants

type EnviromentStruct struct{
	Production string
	Develop		string
}

var EnviromentType = EnviromentStruct {
	Production: "production",
	Develop: "develop",
}