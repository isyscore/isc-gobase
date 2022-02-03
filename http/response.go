package http

type ResponseBase struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DataResponse[T any] struct {
	ResponseBase
	Data T `json:"data"`
}

type PagedData[T any] struct {
	Total         int64 `json:"total"`
	Size          int64 `json:"size"`
	Current       int64 `json:"current"`
	IsSearchCount bool  `json:"isSearchCount"`
	Records       []T   `json:"records"`
}

type PagedResponse[T any] struct {
	ResponseBase
	Data PagedData[T] `json:"data"`
}
