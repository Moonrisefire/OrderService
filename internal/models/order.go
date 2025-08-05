package models

import (
	"time"
)

type Order struct {
	ID          uint      `gorm:"primaryKey"`
	OrderUID    string    `gorm:"uniqueIndex" json:"order_uid" validate:"required,uuid4|alphanum"`
	TrackNumber string    `json:"track_number" validate:"required"`
	Entry       string    `json:"entry" validate:"required"`
	Locale      string    `json:"locale" validate:"required"`
	InternalSig string    `json:"internal_signature"`
	CustomerID  string    `json:"customer_id" validate:"required"`
	DeliverySvc string    `json:"delivery_service" validate:"required"`
	ShardKey    string    `json:"shardkey" validate:"required"`
	SMID        int       `json:"sm_id" validate:"required"`
	DateCreated time.Time `json:"date_created" validate:"required"`
	OofShard    string    `json:"oof_shard" validate:"required"`

	Delivery Delivery `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"delivery" validate:"required,dive"`
	Payment  Payment  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment" validate:"required,dive"`
	Items    []Item   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items" validate:"required,min=1,dive"`
}

type Delivery struct {
	ID      uint   `gorm:"primaryKey"`
	OrderID uint   `gorm:"index"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	ID           uint   `gorm:"primaryKey"`
	OrderID      uint   `gorm:"index"`
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required,len=3"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"required,gte=0"`
	PaymentDT    int64  `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   int    `json:"goods_total" validate:"required,gte=0"`
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ID          uint   `gorm:"primaryKey"`
	OrderID     uint   `gorm:"index"`
	ChrtID      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required,gte=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required,gte=0"`
	NMID        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}
