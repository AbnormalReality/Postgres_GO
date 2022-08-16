package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/AbnormalReality/Postgres_GO/lesson4/pkg//config"
	db "github.com/AbnormalReality/Postgres_GO/lesson4/pkg//database"
	"github.com/AbnormalReality/Postgres_GO/lesson4/pkg//models"
)

type PunishResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func punish(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) PunishResults {
	var queries uint64
	request := *config.Studio

	punisher := func(stopAt time.Time) {
		for {

			studios, err := models.GetStudioByName(ctx, dbpool, request)
			if err != nil {
				log.Fatal(err)
			}
			for _, a := range studios {
				_, err := models.GetGamesByStudioId(ctx, dbpool, a.Id)
				if err != nil {
					log.Fatal(err)
				}
			}

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			punisher(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return PunishResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func main() {
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	dbpool, err := db.InitDBConn(ctx, cnfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	if cnfg.InitDB {
		err = db.InitTables(ctx, dbpool)
		if err != nil {
			log.Fatal(err)
		}
	}

	duration := time.Duration(10 * time.Second)
	threads := 1000
	fmt.Println("start punish")
	res := punish(ctx, duration, threads, dbpool)

	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)
}