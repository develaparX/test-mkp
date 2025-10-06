package transaction

import (
	"database/sql"
	"fmt"
	"time"
)

type TransactionRepo struct {
	DB *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{DB: db}
}

func (r *TransactionRepo) Create(t *Transaction) error {
	query := `
		INSERT INTO transaction (
			id, customer_id, transaction_type, amount,
			transaction_datetime, tax_amount, tax_type,
			payment_status, product_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	var taxType interface{}
	if t.TaxType != "" {
		taxType = t.TaxType
	} else {
		taxType = nil
	}

	_, err := r.DB.Exec(
		query,
		t.ID,
		t.CustomerID,
		t.TransactionType,
		t.Amount,
		t.TransactionDatetime,
		t.TaxAmount,
		taxType,
		t.PaymentStatus,
		t.ProductID,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepo) GetByID(id int64) (*Transaction, error) {
	query := `
		SELECT id, customer_id, transaction_type, amount,
		       transaction_datetime, tax_amount, tax_type,
		       payment_status, product_id
		FROM transaction WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var t Transaction
	var taxType sql.NullString
	var trxTime time.Time

	err := row.Scan(
		&t.ID,
		&t.CustomerID,
		&t.TransactionType,
		&t.Amount,
		&trxTime,
		&t.TaxAmount,
		&taxType,
		&t.PaymentStatus,
		&t.ProductID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan transaction: %w", err)
	}
	t.TransactionDatetime = trxTime
	if taxType.Valid {
		t.TaxType = taxType.String
	}

	return &t, nil
}

func (r *TransactionRepo) GetAll() ([]*Transaction, error) {
	query := `
		SELECT id, customer_id, transaction_type, amount,
		       transaction_datetime, tax_amount, tax_type,
		       payment_status, product_id
		FROM transaction`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		var taxType sql.NullString
		var trxTime time.Time

		err := rows.Scan(
			&t.ID,
			&t.CustomerID,
			&t.TransactionType,
			&t.Amount,
			&trxTime,
			&t.TaxAmount,
			&taxType,
			&t.PaymentStatus,
			&t.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		t.TransactionDatetime = trxTime
		if taxType.Valid {
			t.TaxType = taxType.String
		}
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return transactions, nil
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

func (r *TransactionRepo) GetCustomerActivity(filter CustomerActivityFilter) ([]CustomerActivity, int64, error) {

	baseQuery := `
		WITH ranked AS (
			SELECT
				c.id AS company_id,
				c.name AS company_name,
				cu.id AS customer_id,
				CONCAT(cu.first_name, ' ', cu.last_name) AS full_name,
				COUNT(t.id) AS count_trx
			FROM company c
			INNER JOIN customer cu ON c.id = cu.company
			INNER JOIN transaction t ON cu.id = t.customer_id
			WHERE 1=1`

	var args []interface{}
	argPos := 1

	if filter.CompanyID != nil {
		baseQuery += fmt.Sprintf(" AND c.id = $%d", argPos)
		args = append(args, *filter.CompanyID)
		argPos++
	}

	if filter.MinTrxCount != nil {
		baseQuery += fmt.Sprintf(" AND COUNT(t.id) >= $%d", argPos)
		args = append(args, *filter.MinTrxCount)
		argPos++
	}

	baseQuery += `
			GROUP BY c.id, c.name, cu.id, cu.first_name, cu.last_name
		)
		SELECT
			ROW_NUMBER() OVER (ORDER BY company_id, count_trx DESC, customer_id) AS row_number,
			company_id,
			company_name,
			customer_id,
			full_name,
			count_trx
		FROM ranked`

	countQuery := `
		SELECT COUNT(*) FROM (` + baseQuery + `) AS total`

	var total int64
	err := r.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total rows: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	paginatedQuery := baseQuery + fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.DB.Query(paginatedQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query customer activity: %w", err)
	}
	defer rows.Close()

	var result []CustomerActivity
	for rows.Next() {
		var item CustomerActivity
		err := rows.Scan(
			&item.RowNumber,
			&item.CompanyID,
			&item.CompanyName,
			&item.CustomerID,
			&item.FullName,
			&item.CountTrx,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return result, total, nil
}

func (r *TransactionRepo) GetTransactionsByCustomerAndProduct(customerID, productID int64) ([]Transaction, error) {
	query := `
		SELECT id, customer_id, transaction_type, amount,
		       transaction_datetime, tax_amount, tax_type,
		       payment_status, product_id
		FROM transaction 
		WHERE customer_id = $1 AND product_id = $2
		ORDER BY transaction_datetime DESC`

	rows, err := r.DB.Query(query, customerID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions by customer and product: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		var taxType sql.NullString
		var trxTime time.Time

		err := rows.Scan(
			&t.ID,
			&t.CustomerID,
			&t.TransactionType,
			&t.Amount,
			&trxTime,
			&t.TaxAmount,
			&taxType,
			&t.PaymentStatus,
			&t.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}

		t.TransactionDatetime = trxTime
		if taxType.Valid {
			t.TaxType = taxType.String
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepo) GetTransactionSummaryWithFilter(filter TransactionSummaryFilter) ([]TransactionSummary, int64, error) {

	baseQuery := `
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
		WHERE 1=1`

	var args []interface{}
	argPos := 1

	if filter.CompanyID != nil {
		baseQuery += fmt.Sprintf(" AND c.id = $%d", argPos)
		args = append(args, *filter.CompanyID)
		argPos++
	}

	if filter.ProductID != nil {
		baseQuery += fmt.Sprintf(" AND p.id = $%d", argPos)
		args = append(args, *filter.ProductID)
		argPos++
	}

	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND t.transaction_datetime >= $%d", argPos)
		args = append(args, *filter.StartDate)
		argPos++
	}

	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND t.transaction_datetime <= $%d", argPos)
		args = append(args, *filter.EndDate)
		argPos++
	}

	if filter.MinAmount != nil {
		baseQuery += fmt.Sprintf(" AND t.amount >= $%d", argPos)
		args = append(args, *filter.MinAmount)
		argPos++
	}

	if filter.MaxAmount != nil {
		baseQuery += fmt.Sprintf(" AND t.amount <= $%d", argPos)
		args = append(args, *filter.MaxAmount)
		argPos++
	}

	baseQuery += `
		GROUP BY
			c.id, c.name, p.id, p.product_name, p.service_fee_percentage, p.service_fee
		ORDER BY c.id, p.id`

	countQuery := `SELECT COUNT(*) FROM (` + baseQuery + `) AS total`
	var total int64
	err := r.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total rows: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	paginatedQuery := baseQuery + fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.DB.Query(paginatedQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query transaction summary with filter: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		summaries = append(summaries, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return summaries, total, nil
}
