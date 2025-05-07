package utils

const (
	defaultLimit    = 20
	defaultPage     = 1
	defaultPageSize = 10
	maxLimit        = 200
)

// GetPageAndPageSize validates and returns page size and limit
func GetPageAndPageSize(page, pageSize int) (int, int) {
	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxLimit {
		pageSize = maxLimit
	}
	return page, pageSize
}
