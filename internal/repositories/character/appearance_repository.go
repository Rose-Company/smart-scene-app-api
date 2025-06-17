package character

import (
	"context"
	"smart-scene-app-api/internal/models"
	"smart-scene-app-api/internal/models/character"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppearanceRepository interface {
	Create(ctx context.Context, appearance *character.CharacterAppearance) (*character.CharacterAppearance, error)
	CreateBatch(ctx context.Context, appearances []*character.CharacterAppearance) error
	GetByID(ctx context.Context, id uuid.UUID) (*character.CharacterAppearance, error)
	GetByVideo(ctx context.Context, videoID uuid.UUID) ([]*character.CharacterAppearance, error)
	GetByVideoAndCharacter(ctx context.Context, videoID, characterID uuid.UUID) ([]*character.CharacterAppearance, error)
	List(ctx context.Context, params models.QueryParams, preloadFunc func(*gorm.DB), queryFunc func(*gorm.DB)) ([]*character.CharacterAppearance, error)
	Update(ctx context.Context, id uuid.UUID, appearance *character.CharacterAppearance) (*character.CharacterAppearance, error)
	Delete(ctx context.Context, queryFunc func(*gorm.DB)) error
	DeleteByVideo(ctx context.Context, videoID uuid.UUID) error
	Count(ctx context.Context, queryFunc func(*gorm.DB)) (int64, error)
	GetVideoCharacterSummary(ctx context.Context, videoID uuid.UUID) ([]*character.VideoCharacterSummary, error)
	SearchScenes(ctx context.Context, req *character.CharacterSearchRequest) ([]*character.SceneSegment, error)
}

type appearanceRepository struct {
	db *gorm.DB
}

func NewAppearanceRepository(db *gorm.DB) AppearanceRepository {
	return &appearanceRepository{db: db}
}

func (r *appearanceRepository) Create(ctx context.Context, appearance *character.CharacterAppearance) (*character.CharacterAppearance, error) {
	if err := r.db.WithContext(ctx).Create(appearance).Error; err != nil {
		return nil, err
	}
	return appearance, nil
}

func (r *appearanceRepository) CreateBatch(ctx context.Context, appearances []*character.CharacterAppearance) error {
	return r.db.WithContext(ctx).CreateInBatches(appearances, 100).Error
}

func (r *appearanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*character.CharacterAppearance, error) {
	var appearance character.CharacterAppearance
	if err := r.db.WithContext(ctx).
		Preload("Character").
		Where("id = ?", id).First(&appearance).Error; err != nil {
		return nil, err
	}
	return &appearance, nil
}

func (r *appearanceRepository) GetByVideo(ctx context.Context, videoID uuid.UUID) ([]*character.CharacterAppearance, error) {
	var appearances []*character.CharacterAppearance
	if err := r.db.WithContext(ctx).
		Preload("Character").
		Where("video_id = ?", videoID).
		Order("start_time ASC").
		Find(&appearances).Error; err != nil {
		return nil, err
	}
	return appearances, nil
}

func (r *appearanceRepository) GetByVideoAndCharacter(ctx context.Context, videoID, characterID uuid.UUID) ([]*character.CharacterAppearance, error) {
	var appearances []*character.CharacterAppearance
	if err := r.db.WithContext(ctx).
		Preload("Character").
		Where("video_id = ? AND character_id = ?", videoID, characterID).
		Order("start_time ASC").
		Find(&appearances).Error; err != nil {
		return nil, err
	}
	return appearances, nil
}

func (r *appearanceRepository) List(ctx context.Context, params models.QueryParams, preloadFunc func(*gorm.DB), queryFunc func(*gorm.DB)) ([]*character.CharacterAppearance, error) {
	var appearances []*character.CharacterAppearance

	query := r.db.WithContext(ctx).Model(&character.CharacterAppearance{})

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
	} else {
		query = query.Order("start_time ASC")
	}

	if err := query.Find(&appearances).Error; err != nil {
		return nil, err
	}

	return appearances, nil
}

func (r *appearanceRepository) Update(ctx context.Context, id uuid.UUID, appearance *character.CharacterAppearance) (*character.CharacterAppearance, error) {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(appearance).Error; err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *appearanceRepository) Delete(ctx context.Context, queryFunc func(*gorm.DB)) error {
	query := r.db.WithContext(ctx).Model(&character.CharacterAppearance{})

	if queryFunc != nil {
		queryFunc(query)
	}

	return query.Delete(&character.CharacterAppearance{}).Error
}

func (r *appearanceRepository) DeleteByVideo(ctx context.Context, videoID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("video_id = ?", videoID).Delete(&character.CharacterAppearance{}).Error
}

func (r *appearanceRepository) Count(ctx context.Context, queryFunc func(*gorm.DB)) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&character.CharacterAppearance{})

	if queryFunc != nil {
		queryFunc(query)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *appearanceRepository) GetVideoCharacterSummary(ctx context.Context, videoID uuid.UUID) ([]*character.VideoCharacterSummary, error) {
	var summaries []*character.VideoCharacterSummary

	query := `
		SELECT 
			ca.video_id,
			ca.character_id,
			c.name as character_name,
			c.avatar as character_avatar,
			COUNT(*) as appearance_count,
			SUM(ca.duration) as total_duration,
			MIN(ca.start_time) as first_appearance,
			MAX(ca.end_time) as last_appearance
		FROM character_appearances ca
		JOIN characters c ON ca.character_id = c.id
		WHERE ca.video_id = ?
		GROUP BY ca.video_id, ca.character_id, c.name, c.avatar
		ORDER BY appearance_count DESC
	`

	if err := r.db.WithContext(ctx).Raw(query, videoID).Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

func (r *appearanceRepository) SearchScenes(ctx context.Context, req *character.CharacterSearchRequest) ([]*character.SceneSegment, error) {
	// This would be a complex query to find scenes based on search criteria
	// For now, returning empty slice - would need to implement scene detection logic
	var scenes []*character.SceneSegment
	return scenes, nil
}
