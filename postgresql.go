package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgreManager struct {
	pool *sql.DB
	pipe chan KeyValue
}

var pg *PostgreManager

func ConnPostgre() error {
	psqlconn := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable", "localhost", 5432, "")

	pool, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	pg = &PostgreManager{
		pool: pool,
		pipe: make(chan KeyValue),
	}

	go pg.Run()

	return nil
}

func (p *PostgreManager) Run() {
	defer p.pool.Close()

	for {
		select {
		case msg := <-p.pipe:
			query := fmt.Sprintf("INSERT INTO urls VALUES (%s, %s, CURRENT_DATE)", msg.Key, msg.Value)
			_, err := p.pool.Exec(query)
			if err != nil {
                log.Fatalf("%v", err)
				break
			}
		}
	}
}
