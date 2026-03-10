package url_shortener

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/internal/controller/http/entities"
)

// GetOriginalURLByShortPath принимает json запрос и отправляет на обработку
func (r router) GetOriginalURLByShortPath(c *gin.Context) {
	var urlInput entities.RequestDTOData

	// Проверка на поступивший URL
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		apperor.ErrBadRequest.JsonResponse(c, fmt.Errorf("URL is empty"))
		return
	} else {
		urlInput.URL = shortURL
	}
	// Форматирование поступившей структуры
	pointerURL := urlInput.ToEntity()

	url, err := r.urlUsecase.GetOriginalURLByShortPath(c.Request.Context(), &pointerURL)
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
