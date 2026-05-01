package pagination

import "gorm.io/gorm"

type PageRequest struct {
	Page        int
	PageSize    int
	QueryParams map[string]interface{}
}

type PageResponse struct {
	Total    int64
	Page     int
	PageSize int
	List     interface{}
}

type PageResp[T any] struct {
	Total    int64
	Page     int
	PageSize int
	List     []T
}

func Paginate(pageReq *PageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageReq.Page - 1) * pageReq.PageSize
		return db.Offset(offset).Limit(pageReq.PageSize)
	}
}

func NewPageRequest(page, pageSize int, queryParams map[string]interface{}) *PageRequest {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	return &PageRequest{
		Page:        page,
		PageSize:    pageSize,
		QueryParams: queryParams,
	}
}

func NewPageResponse(total int64, pageReq *PageRequest, list interface{}) *PageResponse {
	return &PageResponse{
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
		List:     list,
	}
}

func NewPageResp[T any](total int64, pageReq *PageRequest, list []T) *PageResp[T] {
	return &PageResp[T]{
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
		List:     list,
	}
}
