package models

type ResponseDatas struct {
	Content string `json:"content"`
	Status  bool   `json:"status"`
}

type CategoryDatas struct {
	Name string `json:"name"`
}
