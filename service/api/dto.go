package api

// Basic Data Transfer Objects common to more components of the API

// ResourceId uniquely identifies a resource in its collection
type ResourceId string

type PaginationInfo struct {
	PageCursor string `json:"pageCursor"`
}

type PageResult[T any] struct {
	HasMoreData bool `json:"hasMoreData"`
	TotalCount  uint `json:"totalCount"`
	PageData    []T  `json:"pageData"`
}
