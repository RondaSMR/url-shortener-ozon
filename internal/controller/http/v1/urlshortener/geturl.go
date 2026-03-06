package urlshortener

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener-ozon/internal/adapters/controller/http/urlapi"
	apperor "url-shortener-ozon/internal/apperror"
)

// GetShortURL принимает json запрос и отправляет на обработку
func (r router) GetShortURL(c *gin.Context) {
	var urlInput urlapi.InOutURL

	shortURL := c.Param("shortURL")
	if shortURL == "" {
		// Если нет в пути, пробуем из JSON (для обратной совместимости)
		if err := c.ShouldBindJSON(&urlInput); err != nil || urlInput.URL == "" {
			apperor.ErrBadRequest.JsonResponse(c, nil)
			return
		}
	} else {
		urlInput.URL = shortURL
	}

	pointerURL, err := urlapi.AdapterHttpURLToEntity(urlInput)
	if err != nil {
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	url, err := r.urlUsecase.GetShortURL(c.Request.Context(), &pointerURL)
	if err != nil {
		if errors.Is(err, apperor.ErrRepoNotFound) {
			apperor.ErrNotFound.JsonResponse(c, err)
			return
		}
		apperor.ErrInternalSystem.JsonResponse(c, err)
		return
	}

	// В случае успеха - делаем редирект на оригинальную ссылку
	c.Redirect(http.StatusMovedPermanently, url.URL)
}
