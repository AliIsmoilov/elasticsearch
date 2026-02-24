package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CategoryI interface {
	Create(context.Context, Category) (*Category, error)
	// Update(context.Context, *UpdateUserReq) (*UserModelResp, error)
	// GetById(context.Context, int64) (*UserModelResp, error)
	// GetByEmail(context.Context, string) (*UserModelResp, error)
	// Delete(context.Context, int64) error
	GetListCategories(ctx context.Context, req GetAllCategoriesReq) (GetAllCategoriesResp, error)
}

type GetAllCategoriesReq struct {
	Limit int32
	Page  int32
	Query string
}

type GetAllCategoriesResp struct {
	Categories []Category
	Count      int64
}

type Category struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
