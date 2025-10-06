package transaction

type TransactionSummary struct {
	ID            int64   `json:"id"`
	CompanyName   string  `json:"company_name"`
	ProductID     int64   `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Amount        float64 `json:"amount"`
	Count         int64   `json:"count"`
	TaxValue      float64 `json:"tax_value"`
	ServiceFeePct bool    `json:"service_fee_percentage"`
	ServiceFee    float64 `json:"service_fee"`
	LastTrxOn     string  `json:"last_trx_on"`
	IDLastTrx     int64   `json:"id_last_trx"`
	FirstTrxOn    string  `json:"first_trx_on"`
	IDFirstTrx    int64   `json:"id_first_trx"`
}
