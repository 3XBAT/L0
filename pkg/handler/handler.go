package handler

import (
	"L0/pkg/service"
	_ "net/http"

	"github.com/gin-gonic/gin"
)

/**/

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*.html")

	order := router.Group("/order")
	{
		order.GET("/:uid", h.getByUID)
	}

	return router

}
