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

var pq *PostgreManager

func ConnPostgre() error {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", "localhost", 5432, "db", "golang", "anon")

	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	pq = &PostgreManager{
		pool: pool,
		pipe: make(chan KeyValue),
	}
	go pq.Run()

	return nil
}

func (p *PostgreManager) Run() {
	defer p.pool.Close()
	defer close(p.pipe)

	for {
		msg := <-p.pipe
		err := p.StoreUrl(msg)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}

func (p *PostgreManager) StoreUrl(msg KeyValue) error {
	var exists bool

	err := p.pool.QueryRow(`SELECT EXISTS (SELECT key FROM urls WHERE key=$1)`, msg.Key).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err := p.pool.Query(`INSERT INTO urls VALUES ($1, $2, CURRENT_DATE)`, msg.Key, msg.Value)
		return err
	}

	return nil
}

func (p *PostgreManager) FindUrl(msg *KeyValue) error {
	err := p.pool.QueryRow(`SELECT value FROM urls WHERE key=$1`, msg.Key).Scan(&msg.Value)

	go func() {
		rdc.pipe <- *msg
	}()

	return err
}
