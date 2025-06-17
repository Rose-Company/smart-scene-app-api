package tag

import (
	"smart-scene-app-api/internal/models"
	tagModels "smart-scene-app-api/internal/models/tag"
	"smart-scene-app-api/server"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	GetTagsHierarchy(req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error)
	GetTagsList(req tagModels.TagFilterRequest) (*models.BaseListResponse, error)
	GetTagsByPosition(req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error)
	GetTagsByCategory(req tagModels.TagFilterRequest) (*models.BaseListResponse, error)
	GetTagUsageStats(limit int, sortBy, sortOrder string) (interface{}, error)
}

type tagService struct {
	sc server.ServerContext
	db *gorm.DB
}

func NewTagService(sc server.ServerContext) Service {
	return &tagService{
		sc: sc,
		db: sc.DB(),
	}
}

func (s *tagService) GetTagsHierarchy(req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error) {
	// TODO: Implement real database query
	// For now, return mock data that matches UI requirements

	response := &tagModels.TagListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    2,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Items: []tagModels.TagHierarchyResponse{
			{
				PositionID:    1,
				PositionTitle: "Video Theme",
				PositionCode:  "video_theme",
				Categories: []tagModels.TagCategoryResponse{
					{
						CategoryID:   1,
						CategoryName: "Gender",
						CategoryCode: "gender",
						Color:        "#007bff",
						FilterType:   "multiple",
						DisplayStyle: "checkbox",
						Tags: []tagModels.TagResponse{
							{TagID: 1, TagName: "Male", TagCode: "male", Color: "#007bff", UsageCount: 150, IsActive: true},
							{TagID: 2, TagName: "Female", TagCode: "female", Color: "#dc3545", UsageCount: 180, IsActive: true},
							{TagID: 3, TagName: "Female & Male", TagCode: "female_male", Color: "#6f42c1", UsageCount: 95, IsActive: true},
						},
					},
				},
			},
			{
				PositionID:    2,
				PositionTitle: "People",
				PositionCode:  "people",
				Categories: []tagModels.TagCategoryResponse{
					{
						CategoryID:   2,
						CategoryName: "Age Range",
						CategoryCode: "age_range",
						Color:        "#28a745",
						FilterType:   "single",
						DisplayStyle: "radio",
						Tags: []tagModels.TagResponse{
							{TagID: 4, TagName: "Child", TagCode: "child", Color: "#28a745", UsageCount: 75, IsActive: true},
							{TagID: 5, TagName: "Teen", TagCode: "teen", Color: "#17a2b8", UsageCount: 120, IsActive: true},
							{TagID: 6, TagName: "Adult", TagCode: "adult", Color: "#ffc107", UsageCount: 200, IsActive: true},
						},
					},
				},
			},
		},
	}

	return response, nil
}

func (s *tagService) GetTagsList(req tagModels.TagFilterRequest) (*models.BaseListResponse, error) {
	// TODO: Implement flat tags list query
	// This would query tags table with joins to categories and positions

	// Mock implementation
	tags := []interface{}{
		map[string]interface{}{
			"id":            1,
			"name":          "Male",
			"code":          "male",
			"color":         "#007bff",
			"usage_count":   150,
			"is_active":     true,
			"category_name": "Gender",
			"category_code": "gender",
		},
		map[string]interface{}{
			"id":            2,
			"name":          "Female",
			"code":          "female",
			"color":         "#dc3545",
			"usage_count":   180,
			"is_active":     true,
			"category_name": "Gender",
			"category_code": "gender",
		},
	}

	response := &models.BaseListResponse{
		Total:    len(tags),
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    tags,
	}

	return response, nil
}

func (s *tagService) GetTagsByPosition(req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error) {
	// TODO: Implement position-specific tags query
	// This would filter by position_code in the request

	var mockItems []tagModels.TagHierarchyResponse

	if req.PositionCode == "people" {
		mockItems = []tagModels.TagHierarchyResponse{
			{
				PositionID:    2,
				PositionTitle: "People",
				PositionCode:  "people",
				Categories: []tagModels.TagCategoryResponse{
					{
						CategoryID:   2,
						CategoryName: "Age Range",
						CategoryCode: "age_range",
						Color:        "#28a745",
						FilterType:   "single",
						DisplayStyle: "radio",
						Tags: []tagModels.TagResponse{
							{TagID: 4, TagName: "Child", TagCode: "child", Color: "#28a745", UsageCount: 75, IsActive: true},
							{TagID: 5, TagName: "Teen", TagCode: "teen", Color: "#17a2b8", UsageCount: 120, IsActive: true},
							{TagID: 6, TagName: "Adult", TagCode: "adult", Color: "#ffc107", UsageCount: 200, IsActive: true},
						},
					},
				},
			},
		}
	} else if req.PositionCode == "video_theme" {
		mockItems = []tagModels.TagHierarchyResponse{
			{
				PositionID:    1,
				PositionTitle: "Video Theme",
				PositionCode:  "video_theme",
				Categories: []tagModels.TagCategoryResponse{
					{
						CategoryID:   1,
						CategoryName: "Gender",
						CategoryCode: "gender",
						Color:        "#007bff",
						FilterType:   "multiple",
						DisplayStyle: "checkbox",
						Tags: []tagModels.TagResponse{
							{TagID: 1, TagName: "Male", TagCode: "male", Color: "#007bff", UsageCount: 150, IsActive: true},
							{TagID: 2, TagName: "Female", TagCode: "female", Color: "#dc3545", UsageCount: 180, IsActive: true},
							{TagID: 3, TagName: "Female & Male", TagCode: "female_male", Color: "#6f42c1", UsageCount: 95, IsActive: true},
						},
					},
				},
			},
		}
	}

	response := &tagModels.TagListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    len(mockItems),
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Items: mockItems,
	}

	return response, nil
}

func (s *tagService) GetTagsByCategory(req tagModels.TagFilterRequest) (*models.BaseListResponse, error) {
	// TODO: Implement category-specific tags query
	// This would filter by category_code in the request

	var mockTags []interface{}

	if req.CategoryCode == "gender" {
		mockTags = []interface{}{
			map[string]interface{}{
				"id":          1,
				"name":        "Male",
				"code":        "male",
				"color":       "#007bff",
				"usage_count": 150,
				"is_active":   true,
			},
			map[string]interface{}{
				"id":          2,
				"name":        "Female",
				"code":        "female",
				"color":       "#dc3545",
				"usage_count": 180,
				"is_active":   true,
			},
			map[string]interface{}{
				"id":          3,
				"name":        "Female & Male",
				"code":        "female_male",
				"color":       "#6f42c1",
				"usage_count": 95,
				"is_active":   true,
			},
		}
	} else if req.CategoryCode == "age_range" {
		mockTags = []interface{}{
			map[string]interface{}{
				"id":          4,
				"name":        "Child",
				"code":        "child",
				"color":       "#28a745",
				"usage_count": 75,
				"is_active":   true,
			},
			map[string]interface{}{
				"id":          5,
				"name":        "Teen",
				"code":        "teen",
				"color":       "#17a2b8",
				"usage_count": 120,
				"is_active":   true,
			},
			map[string]interface{}{
				"id":          6,
				"name":        "Adult",
				"code":        "adult",
				"color":       "#ffc107",
				"usage_count": 200,
				"is_active":   true,
			},
		}
	}

	response := &models.BaseListResponse{
		Total:    len(mockTags),
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    mockTags,
	}

	return response, nil
}

func (s *tagService) GetTagUsageStats(limit int, sortBy, sortOrder string) (interface{}, error) {
	// TODO: Implement real usage statistics query
	// This would aggregate usage counts from video_tags table

	// Mock top tags based on usage
	topTags := []map[string]interface{}{
		{
			"tag_id":      2,
			"tag_name":    "Female",
			"tag_code":    "female",
			"usage_count": 180,
			"percentage":  30.5,
		},
		{
			"tag_id":      1,
			"tag_name":    "Male",
			"tag_code":    "male",
			"usage_count": 150,
			"percentage":  25.4,
		},
		{
			"tag_id":      6,
			"tag_name":    "Adult",
			"tag_code":    "adult",
			"usage_count": 200,
			"percentage":  33.9,
		},
		{
			"tag_id":      5,
			"tag_name":    "Teen",
			"tag_code":    "teen",
			"usage_count": 120,
			"percentage":  20.3,
		},
		{
			"tag_id":      3,
			"tag_name":    "Female & Male",
			"tag_code":    "female_male",
			"usage_count": 95,
			"percentage":  16.1,
		},
	}

	// Apply limit
	if limit > 0 && limit < len(topTags) {
		topTags = topTags[:limit]
	}

	// TODO: Apply sorting based on sortBy and sortOrder parameters

	response := map[string]interface{}{
		"top_tags":     topTags,
		"total_tags":   25,
		"total_usage":  590,
		"generated_at": time.Now().Format(time.RFC3339),
		"sort_by":      sortBy,
		"sort_order":   sortOrder,
		"limit":        limit,
	}

	return response, nil
}

// Private helper methods for real database queries (to be implemented later)

func (s *tagService) buildTagHierarchyQuery(req tagModels.TagFilterRequest) *gorm.DB {
	// TODO: Build complex query with joins
	// SELECT tp.*, tc.*, t.*
	// FROM tag_positions tp
	// LEFT JOIN tag_position_categories tpc ON tp.id = tpc.tag_position_id
	// LEFT JOIN tag_categories tc ON tpc.tag_category_id = tc.id
	// LEFT JOIN tags t ON tc.id = t.category_id
	// WHERE tp.is_active = true AND tc.is_shown = true AND t.is_active = true

	query := s.db.Table("tag_positions tp")

	if req.PositionCode != "" {
		query = query.Where("tp.position = ?", req.PositionCode)
	}

	if req.IsActive != nil {
		query = query.Where("tp.is_active = ?", *req.IsActive)
	}

	return query
}

func (s *tagService) buildTagStatsQuery(sortBy, sortOrder string) *gorm.DB {
	// TODO: Build usage statistics query
	// SELECT t.*, COUNT(vt.video_id) as usage_count
	// FROM tags t
	// LEFT JOIN video_tags vt ON t.id = vt.tag_id
	// GROUP BY t.id
	// ORDER BY usage_count DESC

	query := s.db.Table("tags t").
		Select("t.*, COUNT(vt.video_id) as usage_count").
		Joins("LEFT JOIN video_tags vt ON t.id = vt.tag_id").
		Group("t.id")

	// Apply sorting
	orderClause := "usage_count"
	if sortBy != "" {
		orderClause = sortBy
	}
	if sortOrder == "asc" {
		orderClause += " ASC"
	} else {
		orderClause += " DESC"
	}

	query = query.Order(orderClause)

	return query
}
