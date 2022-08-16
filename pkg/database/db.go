package db

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/AbnormalReality/Postgres_GO/lesson4/pkg/config"
)

func InitDBConn(ctx context.Context, appConfig *config.AppConfig) (dbpool *pgxpool.Pool, err error) {
	// Строка для подключения к базе данных
	url := "postgres://postgres:password@localhost:5432/postgres"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}

	// Pool соединений обязательно ограничивать сверху,
	// так как иначе есть потенциальная опасность превысить лимит соединений с базой.
	cfg.MaxConns = int32(appConfig.MaxConns)
	cfg.MinConns = int32(appConfig.MinConns)

	// HealthCheckPeriod - частота проверки работоспособности
	// соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute

	// MaxConnLifetime - сколько времени будет жить соединение.
	// Так как большого смысла удалять живые соединения нет,
	// можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour

	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	// если запросов не поступало, то соединение закроется.
	cfg.MaxConnIdleTime = 30 * time.Minute

	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second

	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbpool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}

	return
}

func InitTables(ctx context.Context, dbpool *pgxpool.Pool) (err error) {
	query := `
	
create table users (
	id bigint generated always as identity,
	name varchar(200) not null,
	surname varchar(200) not null,
	active boolean default true,
	primary key (id)
);

create table studios (
	id bigint generated always as identity,
	name varchar(200) not null,
    studio_id integer,
	primary key (id)
);

create table gamers (
	id bigint generated always as identity,
	name varchar(200) not null,
	surname varchar(200) not null,
	primary key (id)
);

create table games (
	id bigint generated always as identity,
	title varchar(999) not null,
	studio_id integer,
	constraint fk_studio_id foreign key (studio_id) references studios (id),
	primary key (id)
);

create table gamers_rates (
    date timestamp with time zone default current_timestamp,
    gamer_id integer not null,
    game_id integer not null,
    rate integer not null,
    check (rate > 0 and rate <= 10),
    constraint fk_gamer_id foreign key (gamer_id) references gamers (id),
    constraint fk_game_id foreign key (game_id) references games (id),
    constraint pk_gamer_game primary key(gamer_id, game_id)
);
		
insert into gamers (name, surname) values 
('Sergey','Anikin'), 
('Fox','Mulder'), 
('Dana','Scully');

insert into studios (name, studio_id) values 
('BlueTwelve', 1),('Activision', 2),
('Ubisoft', 3),
('Sony', 4);

insert into games (title, studio_id) values 
('Stray', 1),
('CODMW19', 2),
('RainbowSix', 3),
('GhostofTsushima',4);
		
	
insert into gamers_rates (gamer_id, game_id, rate) values 
(1,1,10), (1,2,10), (1,3,9), (1,4,10), (1,2,8), (1,1,9), (1,2,7),
(2,1,10), (2,2,9), (2,3,10), (2,4,8), (2,1,8), (2,3,9), (2,1,8), 
(3,1,9), (3,2,10), (3,3,9), (3,4,9), (3,4,9), (3,1,9), (3,2,7);
create index concurrently gamers_surname_idx on gamers(surname);
create index concurrently games_title_idx on games using btree (title text_pattern_ops);

	requests := strings.Split(query, ";")

	for _, v := range requests {
		strings.TrimSpace(v)
		if v != "" {
			_, err = dbpool.Exec(ctx, v)
			if err != nil {
				return
			}
		}
	}

	return
}