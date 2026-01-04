package ports

import (
	"database/sql"
	"time"

	"eventure/services/orders/internal/domain/order"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTx(
	tx *sql.Tx,
	orderID string,
	status order.Status,
) error {
	_, err := tx.Exec(
		`INSERT INTO orders (id, status, created_at)
		 VALUES (?, ?, ?)`,
		orderID,
		status,
		time.Now(),
	)
	return err
}

func (r *Repository) UpdateStatusTx(tx *sql.Tx, orderID string, status order.Status) error {
	_, err := tx.Exec(
		`UPDATE orders SET status = ? WHERE id = ?`,
		status,
		orderID,
	)
	return err
}

func (r *Repository) GetStatus(orderID string) (order.Status, error) {
	var status order.Status
	err := r.db.QueryRow(
		`SELECT status FROM orders WHERE id = ?`,
		orderID,
	).Scan(&status)
	return status, err
}
