package response

import "github.com/fajritsaniy/golang-SHM/model/dto"

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type PagedResponse struct {
	Status Status        `json:"status"`
	Data   []interface{} `json:"data,omitempty"`
	Paging dto.Paging    `json:"paging,omitempty"`
}

type FileResponse struct {
	Status   Status `json:"status"`
	FileName string `json:"fileName"`
}
