package video

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	videoModel "smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories/video"
	"smart-scene-app-api/server"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAllVideos(queryParams videoModel.VideoFilterAndPagination) ([]videoModel.Video, error)
	GetVideoByID(id string) (*videoModel.Video, error)
	CreateVideo(video videoModel.Video) (*videoModel.Video, error)
	UpdateVideo(id string, video videoModel.Video) (*videoModel.Video, error)
	DeleteVideo(id string) error
	GetVideosListing(req videoModel.VideoListingRequest) (*videoModel.VideoListResponse, error)
	GetVideoSearchSuggestions(query string, limit int) (*videoModel.VideoSearchSuggestionResponse, error)
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

func (s *videoService) GetAllVideos(queryParams videoModel.VideoFilterAndPagination) ([]videoModel.Video, error) {
	params := models.QueryParams{
		Limit:    queryParams.QueryParams.Limit,
		Offset:   queryParams.QueryParams.Offset,
		Selected: queryParams.QueryParams.Selected,
		QuerySort: models.QuerySort{
			Origin: "created_at desc",
		},
	}
	videos, err := s.videoRepo.List(s.sc.Ctx(), params, func(tx *gorm.DB) {
		if queryParams.Title != "" {
			tx.Where("title ILIKE ?", "%"+queryParams.Title+"%")
		}
		if queryParams.Status != "" {
			tx.Where("status = ?", queryParams.Status)
		}
		if queryParams.CreatedBy != uuid.Nil {
			tx.Where("created_by = ?", queryParams.CreatedBy)
		}
		tx.Preload("CreatedBy").Preload("UpdatedBy")
	}, func(tx *gorm.DB) {
		tx.Order("created_at DESC")
	})

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

	err = s.videoRepo.Delete(s.sc.Ctx(), func(tx *gorm.DB) {
		tx.Where("id = ?", uuidID)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *videoService) GetVideosListing(req videoModel.VideoListingRequest) (*videoModel.VideoListResponse, error) {
	// TODO: Implement real database query with complex search and filters
	// This would involve joins with characters, tags, and other related tables

	// Mock response for now - matches the UI requirements
	mockVideos := []videoModel.VideoListingResponse{
		{
			ID:             uuid.New(),
			Title:          "Sample Video 1",
			ThumbnailURL:   "https://example.com/thumb1.jpg",
			Duration:       120,
			CharacterCount: 5,
			Status:         "completed",
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
			Tags: []videoModel.VideoTagInfo{
				{TagID: 1, TagName: "Tom", TagCode: "tom", TagColor: "#007bff", CategoryID: 1, CategoryName: "Character", Priority: 1},
				{TagID: 2, TagName: "Jerry", TagCode: "jerry", TagColor: "#dc3545", CategoryID: 1, CategoryName: "Character", Priority: 2},
				{TagID: 3, TagName: "Mickey", TagCode: "mickey", TagColor: "#28a745", CategoryID: 1, CategoryName: "Character", Priority: 3},
			},
			VisibleTagsCount: 3,
			TotalTagsCount:   13, // +10 more
			Width:            1920,
			Height:           1080,
			Format:           "mp4",
			FilePath:         "/videos/sample1.mp4",
		},
		{
			ID:             uuid.New(),
			Title:          "Sample Video 2",
			ThumbnailURL:   "https://example.com/thumb2.jpg",
			Duration:       90,
			CharacterCount: 3,
			Status:         "completed",
			CreatedAt:      time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			UpdatedAt:      time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			Tags: []videoModel.VideoTagInfo{
				{TagID: 2, TagName: "Jerry", TagCode: "jerry", TagColor: "#dc3545", CategoryID: 1, CategoryName: "Character", Priority: 1},
				{TagID: 4, TagName: "Spike", TagCode: "spike", TagColor: "#ffc107", CategoryID: 1, CategoryName: "Character", Priority: 2},
			},
			VisibleTagsCount: 2,
			TotalTagsCount:   2,
			Width:            1280,
			Height:           720,
			Format:           "mp4",
			FilePath:         "/videos/sample2.mp4",
		},
	}

	// Apply search filters (mock logic)
	if req.Search != "" {
		// TODO: Filter videos by title search
	}

	if len(req.CharacterNames) > 0 {
		// TODO: Filter videos by character names
	}

	if len(req.TagCodes) > 0 {
		// TODO: Filter videos by tag codes
	}

	// Mock filters for sidebar
	mockFilters := videoModel.VideoListingFilters{
		Statuses: []videoModel.FilterOption{
			{Value: "completed", Label: "Completed", Count: 150},
			{Value: "processing", Label: "Processing", Count: 25},
			{Value: "pending", Label: "Pending", Count: 10},
			{Value: "failed", Label: "Failed", Count: 5},
		},
		DurationRanges: []videoModel.DurationRangeOption{
			{MinDuration: 0, MaxDuration: 30, Label: "0-30s", Count: 50},
			{MinDuration: 30, MaxDuration: 60, Label: "30s-1m", Count: 75},
			{MinDuration: 60, MaxDuration: 300, Label: "1m-5m", Count: 100},
			{MinDuration: 300, MaxDuration: 0, Label: "5m+", Count: 25},
		},
		Tags: []videoModel.TagFilterGroup{
			{
				CategoryID:   1,
				CategoryName: "Gender",
				CategoryCode: "gender",
				FilterType:   "multiple",
				DisplayStyle: "checkbox",
				Tags: []videoModel.TagFilterOption{
					{TagID: 1, TagName: "Male", TagCode: "male", TagColor: "#007bff", Count: 120, IsSelected: false},
					{TagID: 2, TagName: "Female", TagCode: "female", TagColor: "#dc3545", Count: 150, IsSelected: false},
					{TagID: 3, TagName: "Female & Male", TagCode: "female_male", TagColor: "#6f42c1", Count: 80, IsSelected: false},
				},
			},
		},
	}

	response := &videoModel.VideoListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    len(mockVideos),
			Page:     req.Page,
			PageSize: req.PageSize,
			Items:    mockVideos,
		},
		Items:   mockVideos,
		Filters: mockFilters,
	}

	return response, nil
}

func (s *videoService) GetVideoSearchSuggestions(query string, limit int) (*videoModel.VideoSearchSuggestionResponse, error) {
	// TODO: Implement real search suggestions query
	// This would search in video titles and character names

	// Mock response
	response := &videoModel.VideoSearchSuggestionResponse{
		Videos: []videoModel.VideoSearchSuggestion{
			{Type: "video", ID: uuid.New().String(), Title: "Sample Video matching: " + query, ThumbnailURL: "https://example.com/thumb.jpg"},
		},
		Characters: []videoModel.VideoSearchSuggestion{
			{Type: "character", ID: "1", Title: "Tom", Subtitle: "Appears in 25 videos"},
			{Type: "character", ID: "2", Title: "Jerry", Subtitle: "Appears in 30 videos"},
		},
		Total: 3,
	}

	return response, nil
}
