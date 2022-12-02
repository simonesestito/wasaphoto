package api

// Basic Data Transfer Objects common to more components of the API

type PaginationInfo struct {
	PageCursorOrEmpty string `json:"pageCursor" validate:"omitempty,max=80,base64url"`
}

type PageResult[T any] struct {
	NextPageCursor *string `json:"nextPageCursor"`
	PageData       []T     `json:"pageData"`
}
