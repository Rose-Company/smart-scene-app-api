package character

import (
	"net/http"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models/character"

	"github.com/gin-gonic/gin"
)

// GetAllCharacters godoc
// @Summary      Get all characters
// @Description  Retrieve a list of all characters
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.Response{data=[]character.Character}  "List of characters"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/characters [get]
func (h *Handler) GetAllCharacters(c *gin.Context) {
	var queryParams character.CharacterFilterAndPagination
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid query parameters",
			ErrorDetail: err.Error(),
		})
		return
	}

	characters, err := h.service.Character.GetAllCharacters(queryParams)
	if err != nil {
		h.logger.Error("Failed to get characters: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve characters",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Characters retrieved successfully",
		Data:    characters,
	})
}

// GetCharacterByID retrieves a character by its ID
func (h *Handler) GetCharacterByID(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Character ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	character, err := h.service.Character.GetCharacterByID(characterID)
	if err != nil {
		if err == common.ErrCharacterNotFound {
			c.JSON(http.StatusNotFound, common.Response{
				Message:     "Character not found",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get character by ID: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve character",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Character retrieved successfully",
		Data:    character,
	})
}

// CreateCharacter godoc
// @Summary      Create a new character
// @Description  Create a new character with the provided details
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        character  body      character.Character  true  "Character details"
// @Success      201  {object}  common.Response{data=character.Character}  "Character created successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      409  {object}  common.Response  "Character already exists"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/characters [post]
func (h *Handler) CreateCharacter(c *gin.Context) {
	var newCharacter character.Character
	if err := c.ShouldBindJSON(&newCharacter); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid character data",
			ErrorDetail: err.Error(),
		})
		return
	}

	createdCharacter, err := h.service.Character.CreateCharacter(newCharacter)
	if err != nil {
		if err == common.ErrCharacterAlreadyExists {
			c.JSON(http.StatusConflict, common.Response{
				Message:     "Character already exists",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to create character: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to create character",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.Response{
		Message: "Character created successfully",
		Data:    createdCharacter,
	})
}

// UpdateCharacter godoc
// @Summary      Update an existing character
// @Description  Update a character by its ID with the provided details
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Param        character  body      character.Character  true  "Updated character details"
// @Success      200  {object}  common.Response{data=character.Character}  "Character updated successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      404  {object}  common.Response  "Character not found"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/characters/{id} [put]
func (h *Handler) UpdateCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Character ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	var updatedCharacter character.Character
	if err := c.ShouldBindJSON(&updatedCharacter); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid character data",
			ErrorDetail: err.Error(),
		})
		return
	}

	character, err := h.service.Character.UpdateCharacter(characterID, updatedCharacter)
	if err != nil {
		if err == common.ErrCharacterNotFound {
			c.JSON(http.StatusNotFound, common.Response{
				Message:     "Character not found",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to update character: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to update character",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Character updated successfully",
		Data:    character,
	})
}

// DeleteCharacter godoc
// @Summary      Delete a character
// @Description  Delete a character by its ID
// @Tags         characters
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      204  {object}  common.Response  "Character deleted successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      404  {object}  common.Response  "Character not found"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/characters/{id} [delete]
func (h *Handler) DeleteCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Character ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	err := h.service.Character.DeleteCharacter(characterID)
	if err != nil {
		if err == common.ErrCharacterNotFound {
			c.JSON(http.StatusNotFound, common.Response{
				Message:     "Character not found",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to delete character: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to delete character",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Character deleted successfully",
	})
}
