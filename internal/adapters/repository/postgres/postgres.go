package postgres

import (
	"database/sql"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{new(sql.DB)}
}

//
//func (r *UserRepo) GetUser(ctx context.Context, login string) (models.User, error) {
//	return models.User{}, nil
//}
