package video

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	videoModel "smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories/video"
	"smart-scene-app-api/server"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAllVideos() ([]videoModel.Video, error)
	GetVideoByID(id string) (*videoModel.Video, error)
	CreateVideo(video videoModel.Video) (*videoModel.Video, error)
	UpdateVideo(id string, video videoModel.Video) (*videoModel.Video, error)
	DeleteVideo(id string) error
}
type videoService struct {
	sc      server.ServerContext
	videoRepo video.Repository
}
func NewVideoService(sc server.ServerContext) Service {
	return &videoService{
		sc:      sc,
		videoRepo: video.NewRepository(sc.DB()),
	}
}

func (s *videoService) GetAllVideos() ([]videoModel.Video, error) {
	params := models.QueryParams{
		Limit:  10,
		Offset: 0,  
		QuerySort: models.QuerySort{
			Origin: "created_at desc", 
		},
	}
	videos, err := s.videoRepo.List(s.sc.Ctx(), params)
	if err != nil {
		return nil, err
	}
	result := make([]videoModel.Video, len(videos))
	for i, v := range videos {
		if v != nil {
			result[i] = *v
		}
	}
	return result, nil
}

func (s *videoService) GetVideoByID(id string) (*videoModel.Video, error) {
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


	err = s.videoRepo.Delete(s.sc.Ctx(),  func(tx *gorm.DB) {
		tx.Where("id = ?", uuidID)
	})
	if err != nil {
		return err
	}
	return nil
}
