package tag

import (
	"smart-scene-app-api/internal/models"
	tagModels "smart-scene-app-api/internal/models/tag"
	"smart-scene-app-api/server"

	"gorm.io/gorm"
)

type Service interface {
	GetTagsByPosition(req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error)
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

func (s *tagService) GetTagsByPosition(req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error) {
	// TODO: Implement position-specific tags query [Afternoon]
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
