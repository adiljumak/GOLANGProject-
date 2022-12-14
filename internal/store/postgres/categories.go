package postgres

import (
	"Lecture6/internal/models"
	"Lecture6/internal/store"
	"context"
	"github.com/jmoiron/sqlx"
)

func (db *DB) Categories() store.CategoriesRepository {
	if db.categories == nil {
		db.categories = NewCategoriesRepository(db.conn)
	}

	return db.categories
}

type CategoriesRepository struct {
	conn *sqlx.DB
}

func NewCategoriesRepository(conn *sqlx.DB) store.CategoriesRepository {
	return &CategoriesRepository{
		conn: conn,
	}
}

func (c *CategoriesRepository) Create(ctx context.Context, category *models.Category) error {
	_, err := c.conn.Exec("INSERT INTO categories(name) VALUES ($1)", category.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoriesRepository) All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error) {
	basicQuery := "SELECT * FROM categories"

	if filter.Query != nil {
		basicQuery += " WHERE name ilike '%" + *filter.Query + "%'"
	}
	categories := make([]*models.Category, 0)
	if err := c.conn.Select(&categories, basicQuery); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoriesRepository) ByID(ctt context.Context, id int) (*models.Category, error) {
	category := new(models.Category)
	if err := c.conn.Get(category, "SELECT id, name FROM categories WHERE id = $1", id); err != nil {
		return nil, err
	}
	return category, nil
}

func (c *CategoriesRepository) Update(ctx context.Context, category *models.Category) error {
	_, err := c.conn.Exec("UPDATE categories SET name = $1 WHERE id = $2", category.Name, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoriesRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
