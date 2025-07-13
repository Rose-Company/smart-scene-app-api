package video

import (
	"context"
	"smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[video.Video]
	GetVideoTags(ctx context.Context, videoID uuid.UUID) ([]video.VideoTagInfo, error)
	GetVideoTagsMap(ctx context.Context, videoIDs []uuid.UUID) (map[uuid.UUID][]video.VideoTagInfo, error)
}

type repository struct {
	repositories.BaseRepository[video.Video]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := repositories.NewBaseRepository[video.Video](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}

func (r *repository) GetVideoTags(ctx context.Context, videoID uuid.UUID) ([]video.VideoTagInfo, error) {
	var tags []video.VideoTagInfo

	query := `
		SELECT 
			t.id as tag_id,
			t.name as tag_name,
			t.code as tag_code,
			t.color as tag_color,
			tc.id as category_id,
			tc.name as category_name,
			t.sort_order as priority
		FROM video_tags vt
		JOIN tags t ON vt.tag_id = t.id
		JOIN tag_categories tc ON t.category_id = tc.id
		WHERE vt.video_id = ? AND t.is_active = true
		ORDER BY tc.priority ASC, t.sort_order ASC
	`

	err := r.db.WithContext(ctx).Raw(query, videoID).Scan(&tags).Error
	return tags, err
}

func (r *repository) GetVideoTagsMap(ctx context.Context, videoIDs []uuid.UUID) (map[uuid.UUID][]video.VideoTagInfo, error) {
	if len(videoIDs) == 0 {
		return make(map[uuid.UUID][]video.VideoTagInfo), nil
	}

	var results []struct {
		VideoID uuid.UUID `json:"video_id"`
		video.VideoTagInfo
	}

	query := `
		SELECT 
			vt.video_id,
			t.id as tag_id,
			t.name as tag_name,
			t.code as tag_code,
			t.color as tag_color,
			tc.id as category_id,
			tc.name as category_name,
			t.sort_order as priority
		FROM video_tags vt
		JOIN tags t ON vt.tag_id = t.id
		JOIN tag_categories tc ON t.category_id = tc.id
		WHERE vt.video_id IN ? AND t.is_active = true
		ORDER BY vt.video_id, tc.priority ASC, t.sort_order ASC
	`

	err := r.db.WithContext(ctx).Raw(query, videoIDs).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	tagsMap := make(map[uuid.UUID][]video.VideoTagInfo)
	for _, result := range results {
		tagsMap[result.VideoID] = append(tagsMap[result.VideoID], result.VideoTagInfo)
	}

	return tagsMap, nil
}
