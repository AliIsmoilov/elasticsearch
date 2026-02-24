package storage

import (
	esses "elasticsearch/storage/elasticsearch"
	"elasticsearch/storage/postgres"
	"elasticsearch/storage/repo"

	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
)

type StorageI interface {
	Movie() repo.MovieI

	Elastic() repo.ElasticMovieI
}

type storage struct {
	categoryRepo repo.CategoryI
	movieRepo    repo.MovieI
	elasticRepo  repo.ElasticMovieI
}

func New(db *gorm.DB, esClient *elasticsearch.Client) StorageI {
	s := &storage{
		movieRepo: postgres.NewMovie(db),
	}

	if esClient != nil {
		s.elasticRepo = esses.NewMovie(esClient)
	}

	return s
}

func (s *storage) Movie() repo.MovieI {
	return s.movieRepo
}

func (s *storage) Elastic() repo.ElasticMovieI {
	return s.elasticRepo
}
