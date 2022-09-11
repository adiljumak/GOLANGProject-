package store

import (
	"Lecture6/internal/models"
	"context"
)

type Store interface {
	Connect(url string) error
	Close() error

	Categories() CategoriesRepository
	//Phones() GoodsRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error)
	ByID(ctt context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}
type GoodsRepository interface {
	Create(ctx context.Context, phone *models.Good) error
	All(ctx context.Context) ([]*models.Good, error)
	ByID(ctt context.Context, id int) (*models.Good, error)
	Update(ctx context.Context, good *models.Good) error
	Delete(ctx context.Context, id int) error
}
