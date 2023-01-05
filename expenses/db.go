package expenses

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func expService(db *sql.DB) *handler {
	return &handler{db: db}
}

func (h *handler) InitDB() {
	createTB := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err := h.db.Exec(createTB)
	if err != nil {
		log.Fatal("Can't create table", err)
	}
}
