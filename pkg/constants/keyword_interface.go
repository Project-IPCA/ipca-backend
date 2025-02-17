package constants

type ReceiveGetKeyWordPython struct {
	Message string                   `json:"message"`
	Status  string                   `json:"status"`
	Data    GetKeyWordPythonCategory `json:"data"`
}

type GetKeyWordPythonCategory struct {
	Classes       []GetKeyWordCategoryInfo `json:"classes"`
	Functions     []GetKeyWordCategoryInfo `json:"functions"`
	Imports       []GetKeyWordCategoryInfo `json:"imports"`
	Methods       []GetKeyWordCategoryInfo `json:"methods"`
	ReservedWords []GetKeyWordCategoryInfo `json:"reserved_words"`
	Variables     []GetKeyWordCategoryInfo `json:"variables"`
}

type GetKeyWordCategoryInfo struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

type ReceiveGetKeyWordC struct {
	Message string              `json:"message"`
	Status  string              `json:"status"`
	Data    GetKeyWordCCategory `json:"data"`
}

type GetKeyWordCCategory struct {
	Functions    []GetKeyWordCategoryInfo `json:"functions"`
	Includes     []GetKeyWordCategoryInfo `json:"includes"`
	ReverseWords []GetKeyWordCategoryInfo `json:"reserved_words"`
	Variables    []GetKeyWordCategoryInfo `json:"variables"`
}

type CheckKeywordCategoryInfo struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
	Active  bool   `json:"active"`
	Type    string `json:"type"`
}

type PythonCheckKeywordCategory struct {
	Classes       []CheckKeywordCategoryInfo `json:"classes"`
	Functions     []CheckKeywordCategoryInfo `json:"functions"`
	Imports       []CheckKeywordCategoryInfo `json:"imports"`
	Methods       []CheckKeywordCategoryInfo `json:"methods"`
	ReservedWords []CheckKeywordCategoryInfo `json:"reserved_words"`
	Variables     []CheckKeywordCategoryInfo `json:"variables"`
}

type CCheckKeywordCategory struct {
	Functions     []CheckKeywordCategoryInfo `json:"functions"`
	Includes      []CheckKeywordCategoryInfo `json:"includes"`
	ReservedWords []CheckKeywordCategoryInfo `json:"reserved_words"`
	Variables     []CheckKeywordCategoryInfo `json:"variables"`
}

type ResponseCheckKeywordCategoryInfo struct {
	Keyword  string `json:"keyword"`
	Limit    int    `json:"limit"`
	Active   bool   `json:"active"`
	Type     string `json:"type"`
	IsPassed bool   `json:"is_passed"`
}

type PythonCheckKeywordCategoryResponse struct {
	Classes       []ResponseCheckKeywordCategoryInfo `json:"classes"`
	Functions     []ResponseCheckKeywordCategoryInfo `json:"functions"`
	Imports       []ResponseCheckKeywordCategoryInfo `json:"imports"`
	Methods       []ResponseCheckKeywordCategoryInfo `json:"methods"`
	ReservedWords []ResponseCheckKeywordCategoryInfo `json:"reserved_words"`
	Variables     []ResponseCheckKeywordCategoryInfo `json:"variables"`
}

type CCheckKeywordCategoryResponse struct {
	Functions     []ResponseCheckKeywordCategoryInfo `json:"functions"`
	Includes      []ResponseCheckKeywordCategoryInfo `json:"includes"`
	ReservedWords []ResponseCheckKeywordCategoryInfo `json:"reserved_words"`
	Variables     []ResponseCheckKeywordCategoryInfo `json:"variables"`
}

type PythonCheckKeywordResponse struct {
	Status string                             `json:"status"`
	Data   PythonCheckKeywordCategoryResponse `json:"keyword_constraint"`
}

type CCheckKeywordResponse struct {
	Status string                        `json:"status"`
	Data   CCheckKeywordCategoryResponse `json:"keyword_constraint"`
}
