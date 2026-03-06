package urlshortener

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener-ozon/internal/adapters/controller/http/urlapi"
	apperor "url-shortener-ozon/internal/apperror"
	"url-shortener-ozon/pkg/utils"
)

// GetShortURL принимает json запрос и отправляет на обработку
func (r router) GetShortURL(c *gin.Context) {
	var urlInput urlapi.InOutURL

	if err := c.ShouldBindJSON(&urlInput); err != nil {
		apperor.ErrBadRequest.JsonResponse(c, err)
		return
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

	// В случае успеха - возврат оригинальной ссылки
	c.JSON(http.StatusOK, utils.GenerateResponse(nil, urlapi.AdapterEntityToHttpURL(url)))
}
