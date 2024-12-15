package models

type Product struct {
	ID                 uint     `gorm:"primaryKey" json:"id"`
	UserID             uint     `json:"user_id"`
	ProductName        string   `json:"product_name"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `gorm:"type:text[]" json:"product_images"`
	CompressedImages   []string `gorm:"type:text[]" json:"compressed_product_images"`
	ProductPrice       float64  `json:"product_price"`
}
