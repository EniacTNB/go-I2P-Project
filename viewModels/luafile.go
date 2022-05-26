package viewModels

type LuaFile struct {
	FileName    string `json:"file_name"`
	FileContent string `json:"file_content"`
	ExperPoint  string `json:"exper_type"`
	ExperTime   int    `json:"exper_time"`
	Date1       string `json:"date1"`
	Date2       string `json:"date2"`
	IsPublic    bool   `json:"isPublic"`
	StartTime   string `json:"start_time"`
	AddTime     string `json"add_time"`
	Desc        string `json:"desc"`
	User        string `json:"user"`
}
