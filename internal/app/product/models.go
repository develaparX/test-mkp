package product

type Product struct {
	ID                   int64   `json:"id"`
	ProductName          string  `json:"product_name"`
	ServiceFee           float64 `json:"service_fee"`
	ServiceFeePercentage bool    `json:"service_fee_percentage"`
}

var CreateProductReq struct {
	ID                   int64  `json:"id" binding:"required"`
	ProductName          string `json:"product_name" binding:"required,max=100"`
	ServiceFee           string `json:"service_fee" binding:"required"`
	ServiceFeePercentage bool   `json:"service_fee_percentage" binding:"required"`
}

var UpdateProductReq struct {
	ProductName          string `json:"product_name" binding:"required,max=100"`
	ServiceFee           string `json:"service_fee" binding:"required"`
	ServiceFeePercentage bool   `json:"service_fee_percentage" binding:"required"`
}
