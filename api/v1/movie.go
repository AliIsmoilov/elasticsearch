package v1

import (
	"net/http"
	"strconv"

	"elasticsearch/api/models"
	"elasticsearch/storage/repo"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *handlerV1) CreateMovie(ctx *gin.Context) {
	var req models.Movie
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := repo.Movie{
		Rating:       req.Rating,
		Movie:        req.Movie,
		Year:         req.Year,
		Country:      req.Country,
		RatingBall:   req.RatingBall,
		Overview:     req.Overview,
		Director:     req.Director,
		Screenwriter: pq.StringArray(req.Screenwriter),
		Actors:       pq.StringArray(req.Actors),
		UrlLogo:      req.UrlLogo,
	}

	data, err := h.strg.Movie().Create(ctx, m)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := models.Movie{
		Id:           data.Id,
		Rating:       data.Rating,
		Movie:        data.Movie,
		Year:         data.Year,
		Country:      data.Country,
		RatingBall:   data.RatingBall,
		Overview:     data.Overview,
		Director:     data.Director,
		Screenwriter: []string(data.Screenwriter),
		Actors:       []string(data.Actors),
		UrlLogo:      data.UrlLogo,
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (h *handlerV1) GetMovieById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	m, err := h.strg.Movie().GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := models.Movie{
		Id:           m.Id,
		Rating:       m.Rating,
		Movie:        m.Movie,
		Year:         m.Year,
		Country:      m.Country,
		RatingBall:   m.RatingBall,
		Overview:     m.Overview,
		Director:     m.Director,
		Screenwriter: []string(m.Screenwriter),
		Actors:       []string(m.Actors),
		UrlLogo:      m.UrlLogo,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *handlerV1) GetListMovies(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "0")
	page := ctx.DefaultQuery("page", "0")
	query := ctx.DefaultQuery("query", "")

	// parse ints
	// reuse strconv
	limitInt := parseInt32(limit)
	pageInt := parseInt32(page)

	data, err := h.strg.Movie().GetListMovies(ctx, repo.GetAllMoviesReq{
		Limit: limitInt,
		Page:  pageInt,
		Query: query,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := models.GetMoviesListResp{Count: data.Count}
	for _, m := range data.Movies {
		resp.Movies = append(resp.Movies, models.Movie{
			Id:           m.Id,
			Rating:       m.Rating,
			Movie:        m.Movie,
			Year:         m.Year,
			Country:      m.Country,
			RatingBall:   m.RatingBall,
			Overview:     m.Overview,
			Director:     m.Director,
			Screenwriter: []string(m.Screenwriter),
			Actors:       []string(m.Actors),
			UrlLogo:      m.UrlLogo,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

// SearchMoviesES performs a search using Elasticsearch-backed repo if configured.
func (h *handlerV1) SearchMoviesES(ctx *gin.Context) {
	if h.strg.Elastic() == nil {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "elasticsearch not configured"})
		return
	}

	limit := ctx.DefaultQuery("limit", "0")
	page := ctx.DefaultQuery("page", "0")
	query := ctx.DefaultQuery("query", "")

	limitInt := parseInt32(limit)
	pageInt := parseInt32(page)

	data, err := h.strg.Elastic().Search(ctx, repo.GetAllMoviesReq{
		Limit: limitInt,
		Page:  pageInt,
		Query: query,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := models.GetMoviesListResp{Count: data.Count}
	for _, m := range data.Movies {
		resp.Movies = append(resp.Movies, models.Movie{
			Id:           m.Id,
			Rating:       m.Rating,
			Movie:        m.Movie,
			Year:         m.Year,
			Country:      m.Country,
			RatingBall:   m.RatingBall,
			Overview:     m.Overview,
			Director:     m.Director,
			Screenwriter: []string(m.Screenwriter),
			Actors:       []string(m.Actors),
			UrlLogo:      m.UrlLogo,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

func parseInt32(s string) int32 {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return int32(i)
}

func (h *handlerV1) UpdateMovie(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.Movie
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := repo.Movie{
		Id:           id,
		Rating:       req.Rating,
		Movie:        req.Movie,
		Year:         req.Year,
		Country:      req.Country,
		RatingBall:   req.RatingBall,
		Overview:     req.Overview,
		Director:     req.Director,
		Screenwriter: pq.StringArray(req.Screenwriter),
		Actors:       pq.StringArray(req.Actors),
		UrlLogo:      req.UrlLogo,
	}

	data, err := h.strg.Movie().Update(ctx, m)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Movie{
		Id:           data.Id,
		Rating:       data.Rating,
		Movie:        data.Movie,
		Year:         data.Year,
		Country:      data.Country,
		RatingBall:   data.RatingBall,
		Overview:     data.Overview,
		Director:     data.Director,
		Screenwriter: []string(data.Screenwriter),
		Actors:       []string(data.Actors),
		UrlLogo:      data.UrlLogo,
	})
}

func (h *handlerV1) DeleteMovie(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.strg.Movie().Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
