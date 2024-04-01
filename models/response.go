package models

type ResponseNoData struct {
	Status  int    `json:"status_code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

type Response struct {
	Status  int         `json:"status_code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ResponseMultiple struct {
	Status  int         `json:"status_code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Datas   interface{} `json:"datas"`
}
