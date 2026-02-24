package storage

import (
	"elasticsearch/storage/postgres"
	"elasticsearch/storage/repo"

	"gorm.io/gorm"
)

type StorageI interface {
	Movie() repo.MovieI
}

type storage struct {
	categoryRepo repo.CategoryI
	movieRepo    repo.MovieI
}

func New(db *gorm.DB) StorageI {
	return &storage{
		movieRepo: postgres.NewMovie(db),
	}
}

func (s *storage) Movie() repo.MovieI {
	return s.movieRepo
}
