package models

type ResponseDatas struct {
	Content string `json:"content"`
	Status  bool   `json:"status"`
}

type RequestDatas struct {
	Text  string `json:"text"`
	Audio string `json:"voice"`
}
