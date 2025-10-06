
COPY company FROM '/seeds/company.csv' WITH (FORMAT csv, HEADER true, DELIMITER ',');
COPY customer FROM '/seeds/customer.csv' WITH (FORMAT csv, HEADER true, DELIMITER ',');
COPY product FROM '/seeds/product.csv' WITH (FORMAT csv, HEADER true, DELIMITER ',');
COPY transaction FROM '/seeds/transaction.csv' WITH (FORMAT csv, HEADER true, DELIMITER ',');