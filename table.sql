CREATE TABLE IF NOT EXISTS company (
    id BIGINT PRIMARY KEY,
  	name VARCHAR(25) NOT NULL,
	type VARCHAR(25) NOT NULL,
	address VARCHAR(255) NOT NULL,
	city VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS customer (
    id BIGINT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    birth_date DATE,
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20),
    address VARCHAR(255),
    gender VARCHAR(25),
    company BIGINT NOT NULL REFERENCES company(id),
    photo TEXT
);

CREATE TABLE IF NOT EXISTS product (
    id BIGINT PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    service_fee NUMERIC(15, 2) NOT NULL ,
    service_fee_percentage BOOL NOT NULL 
);

CREATE TABLE IF NOT EXISTS transaction (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customer(id),
    transaction_type VARCHAR(20) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL ,
    transaction_datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tax_amount NUMERIC(15, 2) NOT NULL,
    tax_type VARCHAR(50),
    payment_status VARCHAR(20) NOT NULL,
    product_id BIGINT NOT NULL REFERENCES product(id)
);

