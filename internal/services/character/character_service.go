package character

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	characterModel "smart-scene-app-api/internal/models/character"
	characterRepo "smart-scene-app-api/internal/repositories/character"
	"smart-scene-app-api/server"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetAllCharacters(queryParams characterModel.CharacterFilterAndPagination) ([]characterModel.Character, error)
	GetCharacterByID(id string) (*characterModel.Character, error)
	GetCharacterByName(name string) (*characterModel.Character, error)
	CreateCharacter(character characterModel.Character) (*characterModel.Character, error)
	UpdateCharacter(id string, character characterModel.Character) (*characterModel.Character, error)
	DeleteCharacter(id string) error
}

type characterService struct {
	sc            server.ServerContext
	characterRepo characterRepo.Repository
}

func NewCharacterService(sc server.ServerContext) Service {
	return &characterService{
		sc:            sc,
		characterRepo: characterRepo.NewRepository(sc.DB()),
	}
}

func (s *characterService) GetAllCharacters(queryParams characterModel.CharacterFilterAndPagination) ([]characterModel.Character, error) {
	params := models.QueryParams{
		Limit:    queryParams.QueryParams.Limit,
		Offset:   queryParams.QueryParams.Offset,
		Selected: queryParams.QueryParams.Selected,
		QuerySort: models.QuerySort{
			Origin: "created_at desc",
		},
	}

	characters, err := s.characterRepo.List(s.sc.Ctx(), params, func(tx *gorm.DB) {
		// Preload relationships if needed
	}, func(tx *gorm.DB) {
		if queryParams.Name != "" {
			tx.Where("name ILIKE ?", "%"+queryParams.Name+"%")
		}
		if queryParams.IsActive != nil {
			tx.Where("is_active = ?", *queryParams.IsActive)
		}
		if queryParams.CreatedBy != uuid.Nil {
			tx.Where("created_by = ?", queryParams.CreatedBy)
		}
	})

	if err != nil {
		return nil, err
	}

	result := make([]characterModel.Character, len(characters))
	for i, c := range characters {
		if c != nil {
			result[i] = *c
		}
	}
	return result, nil
}

func (s *characterService) GetCharacterByID(id string) (*characterModel.Character, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	character, err := s.characterRepo.GetByID(s.sc.Ctx(), uuidID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrCharacterNotFound
		}
		return nil, err
	}

	return character, nil
}

func (s *characterService) GetCharacterByName(name string) (*characterModel.Character, error) {
	character, err := s.characterRepo.GetByName(s.sc.Ctx(), name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrCharacterNotFound
		}
		return nil, err
	}

	return character, nil
}

func (s *characterService) CreateCharacter(character characterModel.Character) (*characterModel.Character, error) {
	// Check if character with same name already exists
	_, err := s.characterRepo.GetByName(s.sc.Ctx(), character.Name)
	if err == nil {
		return nil, common.ErrCharacterAlreadyExists
	}

	character.ID = uuid.New()
	characterRes, err := s.characterRepo.Create(s.sc.Ctx(), &character)
	if err != nil {
		return nil, err
	}

	return characterRes, nil
}

func (s *characterService) UpdateCharacter(id string, character characterModel.Character) (*characterModel.Character, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	// Check if character exists
	_, err = s.characterRepo.GetByID(s.sc.Ctx(), uuidID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrCharacterNotFound
		}
		return nil, err
	}

	character.ID = uuidID
	updatedCharacter, err := s.characterRepo.Update(s.sc.Ctx(), uuidID, &character)
	if err != nil {
		return nil, err
	}

	return updatedCharacter, nil
}

func (s *characterService) DeleteCharacter(id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return common.ErrInvalidUUID
	}

	// Check if character exists
	_, err = s.characterRepo.GetByID(s.sc.Ctx(), uuidID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.ErrCharacterNotFound
		}
		return err
	}

	err = s.characterRepo.Delete(s.sc.Ctx(), func(tx *gorm.DB) {
		tx.Where("id = ?", uuidID)
	})
	if err != nil {
		return err
	}

	return nil
}
