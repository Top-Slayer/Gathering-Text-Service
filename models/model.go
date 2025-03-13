package models

type ResponseDatas struct {
	Content string `json:"content"`
	Status  bool   `json:"status"`
}

type RequestDatas struct {
	Text  string `json:"text"`
	Audio string `json:"voice"`
}

type ResCheckedDatas struct {
	ID    int64  `json:"id"`
	Text  string `json:"text"`
	Voice string `json:"voice"`
}

type ReqCheckedDatas struct {
	ID     int64  `json:"id"`
	Text   string `json:"text"`
	Status bool   `json:"status"`
}
