package main

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

const dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"

func main() {
	ctx := context.Background()

	// создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	// делаем запрос на создание записи в таблице note
	res, err := con.Exec(ctx, "INSERT INTO notes (title, body) VALUES ($1, $2)", gofakeit.City(), gofakeit.Address().Street)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted %d rows", res.RowsAffected())

	// делаем запрос на получение записей из таблицы note
	rows, err := con.Query(ctx, "SELECT id, title, body, created_at, updated_at FROM note")
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var body string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
	}
}
