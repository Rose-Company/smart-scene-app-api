package character

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		characters := v1.Group("/characters")
		{
			characters.GET("", h.GetAllCharacters)
			characters.POST("", h.CreateCharacter)
			characters.GET("/:id", h.GetCharacterByID)
			characters.PUT("/:id", h.UpdateCharacter)
			characters.DELETE("/:id", h.DeleteCharacter)
		}
	}
}
