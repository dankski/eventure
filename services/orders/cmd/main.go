package main

import (
	"database/sql"
	"log"
	"time"

	"eventure/libs/nats"
	"eventure/services/orders/internal/ports"
	"eventure/services/orders/internal/ports/outbox"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	nc := nats.ConnectNATS()
	defer nc.Close()

	// enable WAL and a busy timeout to reduce "database is locked" errors
	// and limit sql.DB to a single connection to serialize SQLite access
	dsn := "file:orders.db?mode=rwc&_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	// serialize access: sqlite works best with a single open connection
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)

	res, err := db.Exec(string(outbox.SchemaSQL))
	if err != nil {
		log.Fatalf("failed to exec schema %s: %v", outbox.SchemaSQL, err)
	}

	if rows, _ := res.RowsAffected(); rows >= 0 {
		log.Printf("schema applied, rows affected=%d", rows)
	}

	outboxRepo := outbox.NewRepository(db)
	svc := ports.NewOrderService(db, ports.NewRepository(db), outboxRepo)
	svc.StartSagaListeners(nc)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			records, err := outboxRepo.FindUnpublished(10)
			if err != nil {
				log.Println("outbox read error:", err)
				continue
			}

			for _, record := range records {
				if err := nc.Publish(record.EventType, record.Payload); err != nil {
					log.Println("publish failed:", err)
					continue
				}

				if err := outboxRepo.MarkPublished(record.ID); err != nil {
					log.Println("mark published failed:", err)
					continue
				}

				log.Printf("outbox: published event id=%d type=%s", record.ID, record.EventType)
			}
		}
	}()

	for {
		_, err := svc.CreateOrder("item-123", 1)
		if err != nil {
			log.Println("error:", err)
		}
		time.Sleep(5 * time.Second)
	}

}
