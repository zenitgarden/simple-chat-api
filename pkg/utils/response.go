package utils

type Res struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       any         `json:"data"`
	TotalData  int64       `json:"totalData,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}
