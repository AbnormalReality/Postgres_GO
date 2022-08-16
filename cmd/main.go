package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AbnormalReality/Postgres_GO/lesson4/pkg//config"
	"github.com/AbnormalReality/Postgres_GO/lesson4/pkg//database"
	"github.com/AbnormalReality/Postgres_GO/lesson4/pkg//models"
)

func main() {
	err := start()
	if err != nil {
		log.Fatal(err)
	}
}

func start() (err error) {
	cnfg, err := config.NewAppConfig()
	if err != nil {
		return
	}

	ctx := context.Background()
	dbpool, err := db.InitDBConn(ctx, cnfg)
	if err != nil {
		return
	}
	defer dbpool.Close()

	if cnfg.InitDB {
		err = db.InitTables(ctx, dbpool)
		if err != nil {
			return
		}
	}

	request := *config.Studio

	studios, err := models.GetStudioByName(ctx, dbpool, request)
	if err != nil {
		return
	}

	if len(studios) == 0 {
		fmt.Printf("По вашему запросу ничего не найдено\n")
		return
	}

	for _, a := range studios {
		books, err := models.GetGamesByStudioId(ctx, dbpool, a.Id)
		if err != nil {
			return err
		}

		fmt.Printf("По запросу \"%s\" найдены игры:\n", request)
		for _, v := range games {
			fmt.Printf("- %s\n", v.Title)
		}
	}
	return
}