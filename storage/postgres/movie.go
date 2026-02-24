package postgres

import (
	"context"

	"elasticsearch/storage/repo"

	"gorm.io/gorm"
)

type movieRepo struct {
	db *gorm.DB
}

func NewMovie(db *gorm.DB) repo.MovieI {
	return &movieRepo{db: db}
}

func (r *movieRepo) Create(ctx context.Context, m repo.Movie) (*repo.Movie, error) {
	if err := r.db.WithContext(ctx).
		Table("movies").
		Create(&m).
		Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *movieRepo) GetById(ctx context.Context, id int64) (*repo.Movie, error) {
	var m repo.Movie
	if err := r.db.WithContext(ctx).
		Table("movies").
		First(&m, "id = ?", id).
		Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *movieRepo) GetListMovies(ctx context.Context, req repo.GetAllMoviesReq) (repo.GetAllMoviesResp, error) {
	var movies []repo.Movie
	var count int64

	tx := r.db.WithContext(ctx).
		Table("movies").
		Where("deleted_at IS NULL")

	if req.Query != "" {
		tx = tx.Where("movie ILIKE ?", "%"+req.Query+"%")
	}

	if err := tx.Count(&count).Error; err != nil {
		return repo.GetAllMoviesResp{}, err
	}

	if req.Page > 0 && req.Limit > 0 {
		offset := (req.Page - 1) * req.Limit
		tx = tx.Offset(int(offset)).Limit(int(req.Limit))
	} else if req.Limit > 0 {
		tx = tx.Limit(int(req.Limit))
	}

	if err := tx.Find(&movies).Error; err != nil {
		return repo.GetAllMoviesResp{}, err
	}

	return repo.GetAllMoviesResp{Movies: movies, Count: count}, nil
}

func (r *movieRepo) Update(ctx context.Context, m repo.Movie) (*repo.Movie, error) {
	if err := r.db.WithContext(ctx).
		Table("movies").
		Where("id = ?", m.Id).
		Save(&m).
		Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *movieRepo) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).
		Table("movies").
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("now()")).
		Error; err != nil {
		return err
	}
	return nil
}
