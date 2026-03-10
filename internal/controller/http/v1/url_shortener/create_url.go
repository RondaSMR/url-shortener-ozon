package url_shortener

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/controller/http/entities"
	"url-shortener-ozon/pkg/utils"
)

// CreateShortPath принимает json запрос и отправляет на обработку
func (r router) CreateShortPath(c *gin.Context) {
	var urlInput entities.RequestDTOData

	// Парсинг входящего JSON
	if err := c.ShouldBindJSON(&urlInput); err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
	}
	// Форматирование поступившей структуры
	pointerURL := urlInput.ToEntity()

	url, err := r.urlUsecase.CreateShortPath(c.Request.Context(), &pointerURL)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			apperor.ErrValidation.JsonResponse(c, err)
			return
		}
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	// В случае успеха - возврат сокращенной ссылки
	c.JSON(http.StatusCreated, utils.GenerateResponse(nil, entities.FromEntity(url)))
}
