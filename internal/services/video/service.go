package video

import (
	"smart-scene-app-api/internal/models"
	videoModel "smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories/video"
	"smart-scene-app-api/server"
)

type Service interface {
	GetAllVideos() ([]videoModel.Video, error)
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