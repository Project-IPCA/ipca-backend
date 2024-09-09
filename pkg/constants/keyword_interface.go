package constants

type ReceiveGetKeyWordData struct{
	Message string `json:"message"`
	Status string `json:"status"`
	Data GetKeyWordCategory  `json:"data"`
}

type GetKeyWordCategory struct{
	Classes []GetKeyWordCategory  `json:"classes"`
	Functions []GetKeyWordCategory  `json:"functions"`
	Imports []GetKeyWordCategory  `json:"imports"`
	Methods []GetKeyWordCategory  `json:"methods"`
	ReservedWords []GetKeyWordCategory  `json:"reserved_words"`
	Variables []GetKeyWordCategory  `json:"variables"`
}

type GetKeyWordCategoryInfo struct{
	Keyword string  `json:"keyword"`
	Limit int  `json:"limit"`
}

type ReceiveCheckKeywordData struct{
	Status string `json:"status"`
	Data GetKeyWordCategory  `json:"data"`
}

type CheckKeywordCategory struct{
	Classes []CheckKeywordCategoryInfo  `json:"classes"`
	Functions []CheckKeywordCategoryInfo  `json:"functions"`
	Imports []CheckKeywordCategoryInfo  `json:"imports"`
	Methods []CheckKeywordCategoryInfo  `json:"methods"`
	ReservedWords []CheckKeywordCategoryInfo  `json:"reserved_words"`
	Variables []CheckKeywordCategoryInfo  `json:"variables"`
}

type CheckKeywordCategoryInfo struct{
	Keyword string  `json:"keyword"`
	Limit int  `json:"limit"`
    Active bool `json:"active"`
    Type string  `json:"type"`
}

type ResponseCheckKeywordData struct{
	Status string `json:"status"`
	Data ResponseCheckKeywordCategory `json:"keyword_constraint"`
}

type ResponseCheckKeywordCategory struct{
	Classes []ResponseCheckKeywordCategoryInfo  `json:"classes"`
	Functions []ResponseCheckKeywordCategoryInfo  `json:"functions"`
	Imports []ResponseCheckKeywordCategoryInfo  `json:"imports"`
	Methods []ResponseCheckKeywordCategoryInfo  `json:"methods"`
	ReservedWords []ResponseCheckKeywordCategoryInfo  `json:"reserved_words"`
	Variables []ResponseCheckKeywordCategoryInfo  `json:"variables"`
}

type ResponseCheckKeywordCategoryInfo struct{
	Keyword string  `json:"keyword"`
	Limit int  `json:"limit"`
    Active bool `json:"active"`
    Type string  `json:"type"`
	IsPassed bool `json:"is_passed"`
}