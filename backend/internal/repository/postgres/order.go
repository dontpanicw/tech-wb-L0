package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tech-wb-L0/backend/domain"
	"tech-wb-L0/backend/internal/repository/cache"
)

type OrderStorage struct {
	db    *sql.DB
	cache *cache.Cache
}

func (bs *OrderStorage) GetDb() *sql.DB {
	return bs.db
}

func NewOrderStorage(connStr string) (*OrderStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &OrderStorage{db: db,
		cache: cache.NewCache(),
	}, nil
}

func (s *OrderStorage) CreateOrder(ctx context.Context, order *domain.Order) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//delivery
	var deliveryID int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO deliveries (name, phone, zip, city, address, region, email)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
        RETURNING delivery_id`,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		return fmt.Errorf("error insert delivery: %w", err)
	}

	//payment
	var paymentID int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING payment_id`,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee).Scan(&paymentID)
	if err != nil {
		return fmt.Errorf("error insert payment: %w", err)
	}

	//orders
	_, err = tx.ExecContext(ctx, `
        INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
    `, order.OrderUID, order.TrackNumber, order.Entry, deliveryID, paymentID,
		order.Locale, order.InternalSign, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("error insert order: %w", err)
	}

	//items
	for _, item := range order.Items {
		var itemID int64
		err = tx.QueryRowContext(ctx, `
            INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
            RETURNING item_id
        `, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		).Scan(&itemID)
		if err != nil {
			return fmt.Errorf("insert item: %w", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO order_items(order_uid, item_id) VALUES ($1,$2)`,
			order.OrderUID, itemID)
		if err != nil {
			return fmt.Errorf("error insert order_items: %w", err)
		}
	}
	return tx.Commit()
}

func (s *OrderStorage) GetOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	//берем сначала из кеша
	if o, ok := s.cache.Get(orderUID); ok {
		fmt.Println("CACHE HIT:", orderUID)
		return &o, nil
	}
	//достаём order + delivery_id + payment_id
	var o domain.Order
	var deliveryID, paymentID int64

	err := s.db.QueryRowContext(ctx, `
        SELECT o.order_uid, o.track_number, o.entry, 
               o.locale, o.internal_signature, o.customer_id, 
               o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               o.delivery_id, o.payment_id
        FROM orders o
        WHERE o.order_uid = $1
    `, orderUID).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry,
		&o.Locale, &o.InternalSign, &o.CustomerID,
		&o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,
		&deliveryID, &paymentID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // заказа нет
		}
		return nil, fmt.Errorf("select order: %w", err)
	}

	//delivery
	err = s.db.QueryRowContext(ctx, `
        SELECT delivery_id, name, phone, zip, city, address, region, email
        FROM deliveries
        WHERE delivery_id = $1
    `, deliveryID).Scan(
		&o.Delivery.DeliveryID, &o.Delivery.Name, &o.Delivery.Phone,
		&o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address,
		&o.Delivery.Region, &o.Delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("select delivery: %w", err)
	}

	//достаём payment
	err = s.db.QueryRowContext(ctx, `
        SELECT payment_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments
        WHERE payment_id = $1
    `, paymentID).Scan(
		&o.Payment.PaymentID, &o.Payment.Transaction, &o.Payment.RequestID,
		&o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount,
		&o.Payment.PaymentDT, &o.Payment.Bank, &o.Payment.DeliveryCost,
		&o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("select payment: %w", err)
	}

	//достаём items через order_items
	rows, err := s.db.QueryContext(ctx, `
        SELECT i.item_id, i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
        FROM items i
        JOIN order_items oi ON i.item_id = oi.item_id
        WHERE oi.order_uid = $1
    `, orderUID)
	if err != nil {
		return nil, fmt.Errorf("select items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var it domain.Item
		err := rows.Scan(
			&it.ItemID, &it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid,
			&it.Name, &it.Sale, &it.Size, &it.TotalPrice, &it.NmID,
			&it.Brand, &it.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		o.Items = append(o.Items, it)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows items: %w", err)
	}
	s.cache.Set(orderUID, o)
	return &o, nil

}

func (s *OrderStorage) RecoveryCache(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SELECT order_uid FROM orders LIMIT 100`)
	if err != nil {
		return fmt.Errorf("select order_uids: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			log.Printf("Failed to scan order_uid: %v", err)
			continue
		}

		order, err := s.GetOrder(ctx, orderUID)
		if err != nil {
			log.Printf("Failed to get order %s: %v", orderUID, err)
			continue
		}

		s.cache.Set(orderUID, *order)
		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	log.Printf("Cache recovery completed. Restored %d orders", count)
	return nil
}

func (s *OrderStorage) RangeCache() {
	log.Printf("Cache contains %d orders", s.cache.RangeMap())
}
