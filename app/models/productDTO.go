package models

import (
	"time"
)

type Product struct {
	ID              uint            `gorm:"primaryKey"`
	Sku             string          `json:"sku,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	UnitPrice       float32         `json:"unit_price,omitempty"`
	ImageUrl        string          `json:"image_url,omitempty"`
	Active          bool            `json:"active,omitempty"`
	UnitsInStock    int             `json:"units_in_stock,omitempty"`
	DateCreated     time.Time       `gorm:"date_created,omitempty"`
	LastUpdated     time.Time       `gorm:"last_updated,omitempty"`
	CategoryID      uint            `json:"category_id"`
	ProductCategory ProductCategory `gorm:"foreignKey:CategoryID;"`
}

func (Product) TableName() string {
	return "product"
}

type ProductInputDTO struct {
	Sku          *string  `json:"sku,omitempty" form:"Sku"`
	Name         *string  `json:"name,omitempty" form:"Name"`
	Description  *string  `json:"description,omitempty" form:"Description"`
	UnitPrice    *float32 `json:"unit_price,omitempty" form:"UnitPrice"`
	ImageUrl     *string  `json:"image_url,omitempty" form:"ImageUrl"`
	Active       *bool    `json:"active,omitempty" form:"Active"`
	UnitsInStock *int     `json:"units_in_stock,omitempty" form:"UnitsInStock"`
	CategoryID   *uint    `json:"category_id,omitempty" form:"CategoryID"`
}
