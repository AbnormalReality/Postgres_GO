package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type studio struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
}

func Newstudio(name, surname string) *studio {
	return &studio{Name: name}
}

func (a *studio) Add(ctx context.Context, dbpool *pgxpool.Pool) (id int, err error) {
	err = dbpool.QueryRow(ctx, `insert into studios (name) values ($1) returning id`,
		a.Name,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("failed to insert studio: %w", err)
	}
	a.Id = id
	return
}

func (a *studio) Delete(ctx context.Context, dbpool *pgxpool.Pool) (err error) {
	_, err = dbpool.Exec(ctx, `delete from studios where id = $1`, a.Id)
	if err != nil {
		err = fmt.Errorf("failed to delete studio: %w", err)
	}
	return
}

func GetStudioByName(ctx context.Context, dbpool *pgxpool.Pool, name string) (studios []studio, err error) {
	rows, err := dbpool.Query(ctx, `select id, name from studios where name = $1`, name)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a studio
		err = rows.Scan(&a.Id, &a.Name)
		if err != nil {
			err = fmt.Errorf("failed to scan row: %w", err)
			return
		}
		studios = append(studios, a)
	}

	if rows.Err() != nil {
		err = fmt.Errorf("failed to read response: %w", rows.Err())
		return
	}

	return
}