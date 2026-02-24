package api

import (
	v1 "elasticsearch/api/v1"
	"elasticsearch/config"
	"elasticsearch/storage"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Cfg  *config.Config
	Strg storage.StorageI
	Enf  *casbin.Enforcer
}

func New(h *Handler) *gin.Engine {
	engine := gin.Default()

	handlerV1 := v1.New(&v1.HandleV1{
		Cfg:  h.Cfg,
		Strg: h.Strg,
	})

	apiV1 := engine.Group("/v1")

	// movies
	apiV1.POST("/movie", handlerV1.CreateMovie)
	apiV1.GET("/movies", handlerV1.GetListMovies)
	apiV1.GET("/movie/:id", handlerV1.GetMovieById)
	apiV1.PUT("/movie/:id", handlerV1.UpdateMovie)
	apiV1.DELETE("/movie/:id", handlerV1.DeleteMovie)

	return engine
}
