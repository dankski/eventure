package ports

import (
	"database/sql"
	"log"

	"eventure/libs/events"

	status "eventure/services/orders/internal/domain/order"
	"eventure/services/orders/internal/ports/outbox"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type OrderService struct {
	db     *sql.DB
	repo   *Repository
	outbox *outbox.Repository
}

func NewOrderService(db *sql.DB, repo *Repository, outbox *outbox.Repository) *OrderService {
	return &OrderService{db: db, repo: repo, outbox: outbox}
}

func (o *OrderService) CreateOrder(itemID string, qty int) (string, error) {
	orderID := uuid.NewString()

	evt := events.OrderCreated{
		OrderID: orderID,
		ItemID:  itemID,
		Qty:     qty,
	}

	data, _ := events.Marshal(evt)

	tx, err := o.db.Begin()
	if err != nil {
		return "", err
	}

	err = o.repo.CreateTx(tx, orderID, status.StatusNew)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = o.outbox.InsertTx(tx, orderID, "order.created", data)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	log.Printf("OrderService: created order %s", orderID)

	return orderID, nil
}

func (o *OrderService) StartSagaListeners(nc *nats.Conn) {

	nc.Subscribe("inventory.reserved", func(msg *nats.Msg) {
		evt, err := events.UnmarshalInventoryReserved(msg.Data)
		if err != nil {
			log.Println("failed to unmarshal inventory.reserved event:", err)
			return
		}

		tx, err := o.db.Begin()
		if err != nil {
			log.Println("failed to begin transaction:", err)
			return
		}

		err = o.repo.UpdateStatusTx(tx, evt.OrderID, status.StatusReserved)
		if err != nil {
			tx.Rollback()
			log.Println("failed to update order status:", err)
			return
		}

		evtPaymentCharge := events.PaymentCharge{OrderID: evt.OrderID}
		data, _ := events.Marshal(evtPaymentCharge)

		err = o.outbox.InsertTx(tx, evt.OrderID, "payment.charge", data)
		if err != nil {
			tx.Rollback()
			log.Println("failed to insert outbox record:", err)
			return
		}

		tx.Commit()
		log.Printf("OrderService: inventory reserved for order %s", evt.OrderID)
	})

	// inventory failed -> cancel
	nc.Subscribe("inventory.failed", func(msg *nats.Msg) {
		evt, err := events.UnmarshalInventoryFailed(msg.Data)
		if err != nil {
			log.Println("failed to unmarshal inventory.failed event:", err)
			return
		}

		tx, err := o.db.Begin()
		if err != nil {
			log.Println("failed to begin transaction:", err)
			return
		}

		if err := o.repo.UpdateStatusTx(tx, evt.OrderID, status.StatusCancelled); err != nil {
			tx.Rollback()
			log.Println("failed to update order status:", err)
			return
		}

		tx.Commit()
	})

	// payment authorized -> complete order
	nc.Subscribe("payment.authorize", func(msg *nats.Msg) {
		evt, err := events.UnmarshalPaymentAuthorized(msg.Data)
		if err != nil {
			log.Println("failed to unmarshal payment.authorize event:", err)
			return
		}

		tx, err := o.db.Begin()
		if err != nil {
			log.Println("failed to begin transaction:", err)
			return
		}

		if err := o.repo.UpdateStatusTx(tx, evt.OrderID, status.StatusCompleted); err != nil {
			tx.Rollback()
			log.Println("failed to update order status:", err)
			return
		}

		tx.Commit()
		log.Printf("OrderService: payment authorized for order %s", evt.OrderID)
	})

	// payment failed -> cancel
	nc.Subscribe("payment.failed", func(msg *nats.Msg) {
		evt, err := events.UnmarshalPaymentFailed(msg.Data)
		if err != nil {
			log.Println("failed to unmarshal payment.failed event:", err)
			return
		}

		tx, err := o.db.Begin()
		if err != nil {
			log.Println("failed to begin transaction:", err)
			return
		}

		if err := o.repo.UpdateStatusTx(tx, evt.OrderID, status.StatusCancelled); err != nil {
			tx.Rollback()
			log.Println("failed to update order status:", err)
			return
		}
		tx.Commit()
		log.Printf("OrderService: payment failed for order %s", evt.OrderID)
	})
}
