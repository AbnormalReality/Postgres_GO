package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type game struct {
	Id       int    `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	StudioId int    `json:"studio_id" db:"studio_id"`
}

func GetGamesByStudioId(ctx context.Context, dbpool *pgxpool.Pool, studioId int) (games []game, err error) {
	rows, err := dbpool.Query(ctx, `select id, title from games where studio_id = $1`, studioId)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b game
		err = rows.Scan(&b.Id, &b.Title)
		if err != nil {
			err = fmt.Errorf("failed to scan row: %w", err)
			return
		}
		games = append(games, b)
	}

	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		err = fmt.Errorf("failed to read response: %w", rows.Err())
		return
	}

	return
}