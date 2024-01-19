package models

type ProductCategory struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	CategoryName string `json:"category_name"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}
