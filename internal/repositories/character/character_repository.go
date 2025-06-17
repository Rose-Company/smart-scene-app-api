package character

import (
	"context"
	"smart-scene-app-api/internal/models"
	"smart-scene-app-api/internal/models/character"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, character *character.Character) (*character.Character, error)
	GetByID(ctx context.Context, id uuid.UUID) (*character.Character, error)
	GetByName(ctx context.Context, name string) (*character.Character, error)
	List(ctx context.Context, params models.QueryParams, preloadFunc func(*gorm.DB), queryFunc func(*gorm.DB)) ([]*character.Character, error)
	Update(ctx context.Context, id uuid.UUID, character *character.Character) (*character.Character, error)
	Delete(ctx context.Context, queryFunc func(*gorm.DB)) error
	Count(ctx context.Context, queryFunc func(*gorm.DB)) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, character *character.Character) (*character.Character, error) {
	if err := r.db.WithContext(ctx).Create(character).Error; err != nil {
		return nil, err
	}
	return character, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*character.Character, error) {
	var character character.Character
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&character).Error; err != nil {
		return nil, err
	}
	return &character, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*character.Character, error) {
	var character character.Character
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&character).Error; err != nil {
		return nil, err
	}
	return &character, nil
}

func (r *repository) List(ctx context.Context, params models.QueryParams, preloadFunc func(*gorm.DB), queryFunc func(*gorm.DB)) ([]*character.Character, error) {
	var characters []*character.Character

	query := r.db.WithContext(ctx).Model(&character.Character{})

	if queryFunc != nil {
		queryFunc(query)
	}

	if preloadFunc != nil {
		preloadFunc(query)
	}

	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	if params.QuerySort.Origin != "" {
		query = query.Order(params.QuerySort.Origin)
	}

	if err := query.Find(&characters).Error; err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, character *character.Character) (*character.Character, error) {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(character).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *repository) Delete(ctx context.Context, queryFunc func(*gorm.DB)) error {
	query := r.db.WithContext(ctx).Model(&character.Character{})

	if queryFunc != nil {
		queryFunc(query)
	}

	return query.Delete(&character.Character{}).Error
}

func (r *repository) Count(ctx context.Context, queryFunc func(*gorm.DB)) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&character.Character{})

	if queryFunc != nil {
		queryFunc(query)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
