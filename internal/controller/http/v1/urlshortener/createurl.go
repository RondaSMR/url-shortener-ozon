package urlshortener

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"url-shortener-ozon/internal/adapters/controller/http/urlapi"
	"url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/pkg/utils"
)

// CreateShortURL принимает json запрос и отправляет на обработку
func (r router) CreateShortURL(c *gin.Context) {
	var urlInput urlapi.InOutURL

	// Парсинг входящего JSON
	if err := c.ShouldBindJSON(&urlInput); err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
	}
	// Форматирование поступившей структуры
	pointerURL := urlapi.AdapterHttpURLToEntity(urlInput)

	url, err := r.urlUsecase.CreateShortURL(c.Request.Context(), &pointerURL)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			apperor.ErrValidation.JsonResponse(c, err)
			return
		}
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	// В случае успеха - возврат сокращенной ссылки
	c.JSON(http.StatusCreated, utils.GenerateResponse(nil, urlapi.AdapterEntityToHttpURL(url)))
}
