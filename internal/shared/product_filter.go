package shared

// ProductFilter represents filtering criteria for listing products
type ProductFilter struct {
    UserID      uint
    MinPrice    float64
    MaxPrice    float64
    ProductName string
}
