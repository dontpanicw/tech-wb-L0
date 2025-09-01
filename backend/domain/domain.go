package domain

import "time"

type Order struct {
	OrderUID        string    `json:"order_uid"`
	TrackNumber     string    `json:"track_number"`
	Entry           string    `json:"entry"`
	Delivery        Delivery  `json:"delivery"`
	Payment         Payment   `json:"payment"`
	Items           []Item    `json:"items"`
	Locale          string    `json:"locale"`
	InternalSign    string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id"`
	DeliveryService string    `json:"delivery_service"`
	ShardKey        int       `json:"shardkey,string"`
	SmID            int64     `json:"sm_id"`
	DateCreated     time.Time `json:"date_created"`
	OofShard        int       `json:"oof_shard,string"`
}

type Delivery struct {
	DeliveryID int64  `json:"-"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Zip        string `json:"zip"`
	City       string `json:"city"`
	Address    string `json:"address"`
	Region     string `json:"region"`
	Email      string `json:"email"`
}

type Payment struct {
	PaymentID    int64  `json:"-"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ItemID      int64  `json:"-"`
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
