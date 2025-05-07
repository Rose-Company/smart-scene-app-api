package common

type Paging struct {
	Limit *uint `json:"limit" form:"limit"`
	Total int64 `json:"total"`
	Page  *uint `json:"page" form:"page"`
}

var (
	DEFAULT_OFFSET = 0
	DEFAULT_LIMIT  = 10
)

func VerifyPage(page int, limit int) (offset int, finalPage int) {
	if page <= 0 {
		offset = 0
		page = 1
	}
	return (page - 1) * limit, page
}
