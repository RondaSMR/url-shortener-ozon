package url_shortener

import (
	"context"
	"github.com/gin-gonic/gin"
	"url-shortener-ozon/internal/domain/entities"
)

// Usecase описывает логику приложения
type Usecase interface {
	CreateShortPath(ctx context.Context, url *entities.RequestData) (entities.ResponseData, error)
	GetOriginalURLByShortPath(ctx context.Context, url *entities.RequestData) (entities.ResponseData, error)
}

// router хранит зависимость Usecase
type router struct {
	urlUsecase Usecase
}

func Router(
	ginGroup *gin.RouterGroup,
	urlUsecase Usecase,
) {
	r := router{urlUsecase: urlUsecase}

	ginGroup.POST("", r.CreateShortPath)
	ginGroup.GET("/:shortURL", r.GetOriginalURLByShortPath)
}
