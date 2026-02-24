package repo

import (
	"context"
	"time"

	"github.com/lib/pq"
)

type MovieI interface {
	Create(context.Context, Movie) (*Movie, error)
	GetById(context.Context, int64) (*Movie, error)
	GetListMovies(context.Context, GetAllMoviesReq) (GetAllMoviesResp, error)
	Update(context.Context, Movie) (*Movie, error)
	Delete(context.Context, int64) error
}

type ElasticMovieI interface {
	// Search performs a query-based search in Elasticsearch and returns
	// movies and total count matching the request.
	Search(context.Context, GetAllMoviesReq) (GetAllMoviesResp, error)
}

type GetAllMoviesReq struct {
	Limit int32
	Page  int32
	Query string
}

type GetAllMoviesResp struct {
	Movies []Movie
	Count  int64
}

type Movie struct {
	Id           int64          `gorm:"column:id"`
	Rating       int            `gorm:"column:rating"`
	Movie        string         `gorm:"column:movie"`
	Year         int            `gorm:"column:year"`
	Country      string         `gorm:"column:country"`
	RatingBall   float64        `gorm:"column:rating_ball"`
	Overview     string         `gorm:"column:overview"`
	Director     string         `gorm:"column:director"`
	Screenwriter pq.StringArray `gorm:"type:text[];column:screenwriter"`
	Actors       pq.StringArray `gorm:"type:text[];column:actors"`
	UrlLogo      string         `gorm:"column:url_logo"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    *time.Time     `gorm:"column:deleted_at"`
}
