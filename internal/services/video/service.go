package video

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	videoModel "smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories"
	"smart-scene-app-api/internal/repositories/video"
	"smart-scene-app-api/server"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAllVideos(queryParams videoModel.VideoFilterAndPagination) (*videoModel.VideoListResponse, error)
	GetVideoDetail(id string) (*videoModel.Video, error)
	CreateVideo(video videoModel.Video) (*videoModel.Video, error)
	UpdateVideo(id string, video videoModel.Video) (*videoModel.Video, error)
	DeleteVideo(id string) error
}

type videoService struct {
	sc        server.ServerContext
	videoRepo video.Repository
}

func NewVideoService(sc server.ServerContext) Service {
	return &videoService{
		sc:        sc,
		videoRepo: video.NewRepository(sc.DB()),
	}
}

func (s *videoService) GetAllVideos(queryParams videoModel.VideoFilterAndPagination) (*videoModel.VideoListResponse, error) {
	if queryParams.QueryParams.Limit <= 0 {
		queryParams.QueryParams.Limit = 10
	}
	if queryParams.QueryParams.Offset < 0 {
		queryParams.QueryParams.Offset = 0
	}

	var filters []repositories.Clause

	if queryParams.Title != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("title ILIKE ?", "%"+queryParams.Title+"%")
		})
	}
	if queryParams.Status != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("status = ?", queryParams.Status)
		})
	}
	if queryParams.CreatedBy != uuid.Nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("created_by = ?", queryParams.CreatedBy)
		})
	}

	total, err := s.videoRepo.Count(s.sc.Ctx(), models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	page := (queryParams.QueryParams.Offset / queryParams.QueryParams.Limit) + 1
	pageSize := queryParams.QueryParams.Limit

	response := &videoModel.VideoListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    int(total),
			Page:     page,
			PageSize: pageSize,
			Items:    []videoModel.VideoListingResponse{},
		},
	}

	if total == 0 {
		return response, nil
	}

	if queryParams.QueryParams.QuerySort.Origin == "" {
		queryParams.QueryParams.QuerySort.Origin = "created_at.desc"
	}

	filters = append(filters, func(tx *gorm.DB) {
		tx.Preload("CreatedBy").Preload("UpdatedBy")
	})

	videos, err := s.videoRepo.List(s.sc.Ctx(), queryParams.QueryParams, filters...)
	if err != nil {
		return nil, err
	}

	items := make([]videoModel.VideoListingResponse, 0, len(videos))
	for _, v := range videos {
		if v != nil {
			item := videoModel.VideoListingResponse{
				ID:             v.ID,
				Title:          v.Title,
				ThumbnailURL:   v.ThumbnailURL,
				Duration:       v.Duration,
				CharacterCount: v.CharacterCount,
				Status:         v.Status,
				Width:          v.Width,
				Height:         v.Height,
				Format:         v.Format,
				FilePath:       v.FilePath,
				CreatedAt:      v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:      v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
				Tags:           []videoModel.VideoTagInfo{},
			}
			items = append(items, item)
		}
	}

	response.Items = items
	return response, nil
}

// Fix: API Get Video By ID to fit the current models.
func (s *videoService) GetVideoDetail(id string) (*videoModel.Video, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	video, err := s.videoRepo.GetByID(s.sc.Ctx(), uuidID)
	if err != nil {
		return nil, err
	}
	if video == nil {
		return nil, common.ErrVideoNotFound
	}
	return video, nil
}

func (s *videoService) CreateVideo(video videoModel.Video) (*videoModel.Video, error) {
	video.ID = uuid.New()
	videoRes, err := s.videoRepo.Create(s.sc.Ctx(), &video)
	if err != nil {
		return nil, err
	}

	if videoRes == nil {
		return nil, common.ErrVideoNotFound
	}
	return videoRes, nil
}

func (s *videoService) UpdateVideo(id string, video videoModel.Video) (*videoModel.Video, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}
	video.ID = uuidID
	updatedVideo, err := s.videoRepo.Update(s.sc.Ctx(), uuidID, &video)
	if err != nil {
		return nil, err
	}
	return updatedVideo, nil
}

func (s *videoService) DeleteVideo(id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return common.ErrInvalidUUID
	}

	err = s.videoRepo.Delete(s.sc.Ctx(), func(tx *gorm.DB) {
		tx.Where("id = ?", uuidID)
	})
	if err != nil {
		return err
	}
	return nil
}
