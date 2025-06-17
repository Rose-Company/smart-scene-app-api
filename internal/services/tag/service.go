package tag

import (
	"context"
	"smart-scene-app-api/internal/models"
	tagModels "smart-scene-app-api/internal/models/tag"
	tagRepo "smart-scene-app-api/internal/repositories/tag"
	"smart-scene-app-api/server"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	GetTagsByPosition(ctx context.Context, req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error)
}

type tagService struct {
	sc                      server.ServerContext
	db                      *gorm.DB
	tagPositionRepo         *tagRepo.TagPositionRepo
	tagCategoryRepo         *tagRepo.TagCategoryRepo
	tagRepo                 *tagRepo.TagRepo
	tagPositionCategoryRepo *tagRepo.TagPositionCategoryRepo
}

func NewTagService(sc server.ServerContext) Service {
	return &tagService{
		sc:                      sc,
		db:                      sc.DB(),
		tagPositionRepo:         tagRepo.NewTagPositionRepository(sc.DB()),
		tagCategoryRepo:         tagRepo.NewTagCategoryRepository(sc.DB()),
		tagRepo:                 tagRepo.NewTagMainRepository(sc.DB()),
		tagPositionCategoryRepo: tagRepo.NewTagPositionCategoryRepository(sc.DB()),
	}
}

func (s *tagService) GetTagsByPosition(ctx context.Context, req tagModels.TagFilterRequest) (*tagModels.TagListResponse, error) {
	var hierarchyItems []tagModels.TagHierarchyResponse

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	resData := tagModels.TagListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    0,
			Page:     req.Page,
			PageSize: req.PageSize,
			Items:    []tagModels.TagHierarchyResponse{},
		},
	}

	if req.PositionCode == "" {
		return &resData, nil
	}

	position, err := s.tagPositionRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("position = ? AND is_active = ?", req.PositionCode, true)
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &resData, nil
		}
		return nil, err
	}

	positionCategories, err := s.tagPositionCategoryRepo.List(ctx, models.QueryParams{
		QuerySort: models.QuerySort{
			Origin: "tag_position_categories.sort_order.asc",
		},
	}, func(tx *gorm.DB) {
		tx.Select(`
			tag_position_categories.id,
			tag_position_categories.tag_position_id,
			tag_position_categories.tag_category_id,
			tag_position_categories.sort_order,
			tag_position_categories.is_visible,
			tag_position_categories.display_style
		`).
			Preload("TagCategory", func(db *gorm.DB) *gorm.DB {
				if req.CategoryCode != "" {
					return db.Where("code = ?", req.CategoryCode)
				}
				return db.Where("is_shown = ?", true)
			}).
			Where("tag_position_id = ? AND is_visible = ?", position.ID, true)
	})
	if err != nil {
		return nil, err
	}

	if len(positionCategories) == 0 {
		return &resData, nil
	}

	categories := make([]tagModels.TagCategoryResponse, 0)

	for _, pc := range positionCategories {
		if pc.TagCategory.ID == 0 {
			continue
		}

		if req.CategoryCode != "" && pc.TagCategory.Code != req.CategoryCode {
			continue
		}

		tags, err := s.tagRepo.List(ctx, models.QueryParams{
			QuerySort: models.QuerySort{
				Origin: "sort_order.asc,name.asc",
			},
		}, func(tx *gorm.DB) {
			conditions := []string{"category_id = ?"}
			values := []interface{}{pc.TagCategory.ID}

			if req.IsActive != nil {
				conditions = append(conditions, "is_active = ?")
				values = append(values, *req.IsActive)
			} else {
				conditions = append(conditions, "is_active = ?")
				values = append(values, true)
			}

			if req.IsSystemTag != nil {
				conditions = append(conditions, "is_system_tag = ?")
				values = append(values, *req.IsSystemTag)
			}

			tx.Where(strings.Join(conditions, " AND "), values...)
		})

		if err != nil || len(tags) == 0 {
			continue
		}

		tagResponses := make([]tagModels.TagResponse, 0)
		for _, tag := range tags {
			color := tag.Color
			if color == "" {
				color = pc.TagCategory.Color
			}

			tagResponses = append(tagResponses, tagModels.TagResponse{
				TagID:      tag.ID,
				TagName:    tag.Name,
				TagCode:    tag.Code,
				Color:      color,
				UsageCount: tag.UsageCount,
				IsActive:   tag.IsActive,
			})
		}

		categories = append(categories, tagModels.TagCategoryResponse{
			CategoryID:   pc.TagCategory.ID,
			CategoryName: pc.TagCategory.Name,
			CategoryCode: pc.TagCategory.Code,
			Color:        pc.TagCategory.Color,
			FilterType:   pc.TagCategory.FilterType,
			DisplayStyle: pc.DisplayStyle,
			Tags:         tagResponses,
		})
	}

	if len(categories) > 0 {
		hierarchyItems = append(hierarchyItems, tagModels.TagHierarchyResponse{
			PositionID:    position.ID,
			PositionTitle: position.Title,
			PositionCode:  position.Position,
			Categories:    categories,
		})
	}

	resData.Items = hierarchyItems
	resData.Total = len(hierarchyItems)

	return &resData, nil
}
