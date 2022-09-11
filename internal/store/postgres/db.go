package postgres

import (
	"Lecture6/internal/store"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB

	goods      store.GoodsRepository
	categories store.CategoriesRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(url string) error {
	s := pgx.Conn{}
	s = s
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	//if err := conn.Ping(); err != nil {
	//	return err
	//}

	db.conn = conn
	return nil
}
func (db *DB) Close() error {
	return db.conn.Close()
}
