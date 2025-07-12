package video

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	videoModel "smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories"
	"smart-scene-app-api/internal/repositories/video"
	"smart-scene-app-api/server"

	"strings"

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
	queryParams.VerifyPaging()

	limit := queryParams.PageSize
	offset := (queryParams.Page - 1) * queryParams.PageSize

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

	if len(queryParams.TagIDs) > 0 || len(queryParams.TagCodes) > 0 {
		filters = append(filters, func(tx *gorm.DB) {
			subQuery := s.sc.DB().
				Table("video_tags vt").
				Select("vt.video_id").
				Joins("JOIN tags t ON vt.tag_id = t.id").
				Where("t.is_active = ?", true)

			if len(queryParams.TagIDs) > 0 {
				subQuery = subQuery.Where("vt.tag_id IN ?", queryParams.TagIDs)
			}
			if len(queryParams.TagCodes) > 0 {
				var codes []string
				for _, code := range queryParams.TagCodes {
					splitCodes := strings.Split(code, ",")
					codes = append(codes, splitCodes...)
				}
				subQuery = subQuery.Where("t.code IN ?", codes)
			}

			tx.Where("id IN (?)", subQuery)
		})
	}

	total, err := s.videoRepo.Count(s.sc.Ctx(), models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	response := &videoModel.VideoListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    int(total),
			Page:     queryParams.Page,
			PageSize: queryParams.PageSize,
			Items:    []videoModel.VideoListingResponse{},
		},
	}

	if total == 0 {
		return response, nil
	}

	sort := queryParams.Sort
	if sort == "" {
		sort = "created_at.desc"
	}

	repoQueryParams := models.QueryParams{
		Limit:  limit,
		Offset: offset,
		QuerySort: models.QuerySort{
			Origin: sort,
		},
	}

	videos, err := s.videoRepo.List(s.sc.Ctx(), repoQueryParams, filters...)
	if err != nil {
		return nil, err
	}

	videoIDs := make([]uuid.UUID, 0, len(videos))
	for _, v := range videos {
		if v != nil {
			videoIDs = append(videoIDs, v.ID)
		}
	}

	tagsMap, err := s.videoRepo.GetVideoTagsMap(s.sc.Ctx(), videoIDs)
	if err != nil {
		return nil, err
	}

	items := make([]videoModel.VideoListingResponse, 0, len(videos))
	for _, v := range videos {
		if v != nil {
			tags := tagsMap[v.ID]
			if tags == nil {
				tags = []videoModel.VideoTagInfo{}
			}

			item := videoModel.VideoListingResponse{
				ID:               v.ID,
				Title:            v.Title,
				ThumbnailURL:     v.ThumbnailURL,
				Duration:         v.Duration,
				CharacterCount:   v.CharacterCount,
				Status:           v.Status,
				FilePath:         v.FilePath,
				CreatedAt:        v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:        v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
				Tags:             tags,
				VisibleTagsCount: len(tags),
				TotalTagsCount:   len(tags),
				Metadata:         v.Metadata,
			}
			items = append(items, item)
		}
	}

	response.Items = items
	return response, nil
}

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
