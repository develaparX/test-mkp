package transaction

import (
	"database/sql"
	"fmt"
)

type TransactionRepo struct {
	DB *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{DB: db}
}

func (r *TransactionRepo) GetTransactionSummary() ([]TransactionSummary, error) {
	query := `
SELECT
    c.id,
    c.name AS company_name,
    p.id AS product_id,
    p.product_name,
    SUM(t.amount) AS amount,
    COUNT(t.id) AS count,
    SUM(t.tax_amount) AS tax_value,
    p.service_fee_percentage,
    p.service_fee,
    MAX(t.transaction_datetime) AS last_trx_on,
    (
        SELECT t2.id
        FROM transaction t2
        INNER JOIN customer cu2 ON t2.customer_id = cu2.id
        WHERE cu2.company = c.id AND t2.product_id = p.id
        ORDER BY t2.transaction_datetime DESC, t2.id DESC
        LIMIT 1
    ) AS id_last_trx,
    MIN(t.transaction_datetime) AS first_trx_on,
    (
        SELECT t3.id
        FROM transaction t3
        INNER JOIN customer cu3 ON t3.customer_id = cu3.id
        WHERE cu3.company = c.id AND t3.product_id = p.id
        ORDER BY t3.transaction_datetime ASC, t3.id ASC
        LIMIT 1
    ) AS id_first_trx
FROM
    company c
    INNER JOIN customer cu ON c.id = cu.company
    INNER JOIN transaction t ON cu.id = t.customer_id
    INNER JOIN product p ON t.product_id = p.id
GROUP BY
    c.id, c.name, p.id, p.product_name, p.service_fee_percentage, p.service_fee
ORDER BY c.id, p.id;
`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query transaction summary: %w", err)
	}
	defer rows.Close()

	var summaries []TransactionSummary
	for rows.Next() {
		var s TransactionSummary
		err := rows.Scan(
			&s.ID,
			&s.CompanyName,
			&s.ProductID,
			&s.ProductName,
			&s.Amount,
			&s.Count,
			&s.TaxValue,
			&s.ServiceFeePct,
			&s.ServiceFee,
			&s.LastTrxOn,
			&s.IDLastTrx,
			&s.FirstTrxOn,
			&s.IDFirstTrx,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		summaries = append(summaries, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return summaries, nil
}
