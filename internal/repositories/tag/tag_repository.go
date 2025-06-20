package tag

import (
	tagModels "smart-scene-app-api/internal/models/tag"
	"smart-scene-app-api/internal/repositories"
	"smart-scene-app-api/server"

	"gorm.io/gorm"
)

type Repository interface {
}

type repository struct {
	sc                      server.ServerContext
	db                      *gorm.DB
	tagPositionRepo         TagPositionRepo
	tagCategoryRepo         TagCategoryRepo
	tagRepo                 TagRepo
	tagPositionCategoryRepo TagPositionCategoryRepo
}

type TagPositionRepo struct {
	db *gorm.DB
	repositories.BaseRepository[tagModels.TagPosition]
}

type TagCategoryRepo struct {
	db *gorm.DB
	repositories.BaseRepository[tagModels.TagCategory]
}

type TagRepo struct {
	db *gorm.DB
	repositories.BaseRepository[tagModels.Tag]
}

type TagPositionCategoryRepo struct {
	db *gorm.DB
	repositories.BaseRepository[tagModels.TagPositionCategory]
}

func NewTagRepository(sc server.ServerContext) Repository {
	return &repository{
		sc:                      sc,
		db:                      sc.DB(),
		tagPositionRepo:         *NewTagPositionRepository(sc.DB()),
		tagCategoryRepo:         *NewTagCategoryRepository(sc.DB()),
		tagRepo:                 *NewTagMainRepository(sc.DB()),
		tagPositionCategoryRepo: *NewTagPositionCategoryRepository(sc.DB()),
	}
}

func NewTagPositionRepository(db *gorm.DB) *TagPositionRepo {
	baseRepo := repositories.NewBaseRepository[tagModels.TagPosition](db)
	return &TagPositionRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func NewTagCategoryRepository(db *gorm.DB) *TagCategoryRepo {
	baseRepo := repositories.NewBaseRepository[tagModels.TagCategory](db)
	return &TagCategoryRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func NewTagMainRepository(db *gorm.DB) *TagRepo {
	baseRepo := repositories.NewBaseRepository[tagModels.Tag](db)
	return &TagRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func NewTagPositionCategoryRepository(db *gorm.DB) *TagPositionCategoryRepo {
	baseRepo := repositories.NewBaseRepository[tagModels.TagPositionCategory](db)
	return &TagPositionCategoryRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
