package urlshortener

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener-ozon/internal/adapters/controller/http/urlapi"
	apperor "url-shortener-ozon/internal/apperror"
)

// GetShortURL принимает json запрос и отправляет на обработку
func (r router) GetShortURL(c *gin.Context) {
	var urlInput urlapi.InOutURL

	// Проверка на поступивший URL
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		apperor.ErrBadRequest.JsonResponse(c, fmt.Errorf("URL is empty"))
		return
	} else {
		urlInput.URL = shortURL
	}

	pointerURL := urlapi.AdapterHttpURLToEntity(urlInput)

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
