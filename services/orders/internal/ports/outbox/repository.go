package outbox

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertTx(
	tx *sql.Tx,
	aggregateID string,
	eventType string,
	payload []byte,
) error {
	_, err := tx.Exec(
		`INSERT INTO outbox (aggregate_id, event_type, payload) VALUES (?, ?, ?)`,
		aggregateID,
		eventType,
		payload,
	)
	return err
}

func (r *Repository) FindUnpublished(limit int) ([]Record, error) {
	rows, err := r.db.Query(
		`SELECT id, aggregate_id, event_type, payload FROM outbox WHERE published_at IS NULL ORDER BY id LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.AggregateID, &rec.EventType, &rec.Payload); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *Repository) MarkPublished(id int64) error {
	_, err := r.db.Exec(
		`UPDATE outbox SET published_at = ? WHERE id = ?`,
		time.Now(),
		id,
	)
	return err
}
