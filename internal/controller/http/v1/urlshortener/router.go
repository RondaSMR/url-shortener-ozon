package urlshortener

import (
	"context"
	"github.com/gin-gonic/gin"
	"url-shortener-ozon/internal/domain/entities"
)

// Usecase описывает логику приложения
type Usecase interface {
	CreateShortURL(ctx context.Context, url *entities.InOutURL) (entities.InOutURL, error)
	GetShortURL(ctx context.Context, url *entities.InOutURL) (entities.InOutURL, error)
}

// router хранит зависимость Usecase
type router struct {
	urlUsecase Usecase
}

func Router(
	ginGroup *gin.RouterGroup,
	urlUsecase Usecase,
	user string,
	pass string,
) {
	r := router{urlUsecase: urlUsecase}

	ginGroup.Use(gin.BasicAuth(gin.Accounts{
		user: pass,
	}))

	ginGroup.POST("", r.CreateShortURL)
	ginGroup.GET("", r.GetShortURL)
}
