DROP DATABASE IF EXISTS fresh_products;
CREATE DATABASE fresh_products;
USE fresh_products;


-- Sprint 1, requirement 1
CREATE TABLE sellers(
    id INT PRIMARY KEY AUTO_INCREMENT,
    cid INT,
    company_name VARCHAR(255),
    address VARCHAR(255),
    telephone VARCHAR(255),
    locality_id INT
);

-- Sprint 1, requirement 2
CREATE TABLE warehouses(
    id INT PRIMARY KEY AUTO_INCREMENT,
    address VARCHAR(255),
    telephone VARCHAR(255),
    warehouse_code VARCHAR(255),
    locality_id INT
);

-- Sprint 1, requirement 3
CREATE TABLE sections(
    id INT PRIMARY KEY AUTO_INCREMENT,
    -- section_number é int, não varchar
    section_number INT,
    current_capacity INT,
    maximum_capacity INT,
    minimum_capacity INT,
    current_temperature DECIMAL(19,2),
    minimum_temperature DECIMAL(19,2),
    product_type_id INT,
    warehouse_id INT
);

-- Sprint 1, requirement 4
CREATE TABLE products(
    id INT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(255),
    expiration_rate DECIMAL(19,2),
    freezing_rate DECIMAL(19,2),
    height DECIMAL(19,2),
    length DECIMAL(19,2),
    net_weight DECIMAL(19,2),
    product_code VARCHAR(255),
    recommended_freezing_temperature DECIMAL(19,2),
    width DECIMAL(19,2),
    product_type_id INT,
    seller_id INT
);
CREATE TABLE product_types(
    id INT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(255)
);

-- Sprint 1, requirement 5
CREATE TABLE employees(
    id INT PRIMARY KEY AUTO_INCREMENT,
    id_card_number VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    warehouse_id INT
);

-- Sprint 1, requirement 6
CREATE TABLE buyers(
    id INT PRIMARY KEY AUTO_INCREMENT,
    id_card_number VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255)
);


-- Sprint 2, requirement 1
CREATE TABLE localities(
    id INT PRIMARY KEY NOT NULL,
    locality_name VARCHAR(255),
    province_id INT
);

CREATE TABLE provinces(
    id INT PRIMARY KEY NOT NULL,
    province_name VARCHAR(255),
    id_country_fk INT
);
CREATE TABLE countries(
    id INT PRIMARY KEY NOT NULL,
    country_name VARCHAR(255)
);


-- Sprint 2, requirement 2
CREATE TABLE carriers(
    id INT PRIMARY KEY AUTO_INCREMENT,
    cid VARCHAR(255),
    company_name VARCHAR(255),
    address VARCHAR(255),
    telephone VARCHAR(255),
    locality_id INT
);

-- Sprint 2, requirement 3
CREATE TABLE product_batches(
    id INT PRIMARY KEY AUTO_INCREMENT,
    batch_number VARCHAR(255),
    current_quantity INT,
    current_temperature DECIMAL(19,2),
    due_date DATETIME(6),
    initial_quantity INT,
    manufacturing_date DATETIME(6),
    manufacturing_hour INT(2),
    minimum_temperature DECIMAL(19,2),
    product_id INT,
    section_id INT
);

-- Sprint 2, requirement 4
CREATE TABLE product_records(
    id INT PRIMARY KEY AUTO_INCREMENT,
    last_update_date DATETIME(6),
    purchase_price DECIMAL(19,2),
    sale_price DECIMAL(19,2),
    product_id INT
);

-- Sprint 2, requirement 5
CREATE TABLE inbound_orders(
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_date DATETIME(6),
    order_number VARCHAR(255),
    employee_id INT,
    product_batch_id INT,
    warehouse_id INT
);

-- Sprint 2, requirement 6
CREATE TABLE purchase_orders(
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_number VARCHAR(255),
    order_date DATETIME(6),
    tracking_code VARCHAR(255),
    buyer_id INT,
    carrier_id INT,
    order_status_id INT,
    warehouse_id INT
);
CREATE TABLE order_status(
    id INT PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(255)
);


-- Sprint 1 constraints
-- R1
ALTER TABLE sellers ADD FOREIGN KEY (locality_id) REFERENCES localities(id);
-- R2
ALTER TABLE warehouses ADD FOREIGN KEY (locality_id) REFERENCES localities(id);
-- R3
ALTER TABLE sections ADD FOREIGN KEY (warehouse_id) REFERENCES warehouses(id);
ALTER TABLE sections ADD FOREIGN KEY (product_type_id) REFERENCES product_types(id);
-- R4
ALTER TABLE products ADD FOREIGN KEY (product_type_id) REFERENCES product_types(id);
ALTER TABLE products ADD FOREIGN KEY (seller_id) REFERENCES sellers(id);
-- R5
ALTER TABLE employees ADD FOREIGN KEY (warehouse_id) REFERENCES warehouses(id);
-- R6

-- Sprint 2 constraints
-- R1
ALTER TABLE localities ADD FOREIGN KEY (province_id) REFERENCES provinces(id);
-- R2
ALTER TABLE carriers ADD FOREIGN KEY (locality_id) REFERENCES localities(id);
-- R3
ALTER TABLE product_batches ADD FOREIGN KEY (product_id) REFERENCES products(id);
ALTER TABLE product_batches ADD FOREIGN KEY (section_id) REFERENCES sections(id);
-- R4
ALTER TABLE product_records ADD FOREIGN KEY (product_id) REFERENCES products(id);
-- R5
ALTER TABLE inbound_orders ADD FOREIGN KEY (employee_id) REFERENCES employees(id);
ALTER TABLE inbound_orders ADD FOREIGN KEY (product_batch_id) REFERENCES product_batches(id);
ALTER TABLE inbound_orders ADD FOREIGN KEY (warehouse_id) REFERENCES warehouses(id);
-- R6
ALTER TABLE purchase_orders ADD FOREIGN KEY (warehouse_id) REFERENCES warehouses(id);
ALTER TABLE purchase_orders ADD FOREIGN KEY (buyer_id) REFERENCES buyers(id);
ALTER TABLE purchase_orders ADD FOREIGN KEY (carrier_id) REFERENCES carriers(id);
ALTER TABLE purchase_orders ADD FOREIGN KEY (order_status_id) REFERENCES order_status(id);





-- Insert sample countries
INSERT INTO countries (id, country_name) VALUES
(1, 'USA'),
(2, 'Canada');

-- Insert sample provinces
INSERT INTO provinces (id, province_name, id_country_fk) VALUES
(1, 'California', 1),
(2, 'Ontario', 2);

-- Insert sample localities
INSERT INTO localities (id, locality_name, province_id) VALUES
(1, 'Los Angeles', 1),
(2, 'Toronto', 2),
(3, 'Vancouver', 2);

-- Insert sample carriers
INSERT INTO carriers (id, cid, company_name, address, telephone, locality_id) VALUES
(1, 'CARRIER001', 'Fast Logistics', '123 Main St, LA', '555-1234', 1),
(2, 'CARRIER002', 'Speedy Delivery', '456 Oak St, Toronto', '555-5678', 2);

-- Insert sample order statuses
INSERT INTO order_status (id, description) VALUES
(1, 'Pending'),
(2, 'Shipped'),
(3, 'Delivered');

-- Insert sample sellers
INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES
(100, 'Fresh Foods Co.', '789 Fruit Rd, LA', '555-9876', 1),
(200, 'Organic Produce Ltd.', '101 Veggie Blvd, Toronto', '555-2345', 2);

-- Insert sample warehouses
INSERT INTO warehouses (address, telephone, warehouse_code, locality_id) VALUES
('1234 Cold Storage St, LA', '555-3456', 'WH001', 1),
('5678 Cool Goods Ave, Toronto', '555-7890', 'WH002', 2);

-- Insert sample product types
INSERT INTO product_types (description) VALUES
('Vegetables'),
('Fruits'),
('Dairy'),
('Meat');

-- Insert sample products
INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES
('Fresh Apples', 0.1, 0.2, 4.5, 7.5, 1.2, 'PA001', -2.5, 5.0, 2, 1),
('Organic Carrots', 0.2, 0.1, 6.0, 5.0, 0.8, 'CA001', -3.0, 4.0, 1, 2),
('Skimmed Milk', 0.05, 0.03, 10.0, 15.0, 1.0, 'MI001', -4.0, 8.0, 3, 1),
('Chicken Breasts', 0.3, 0.4, 8.0, 10.0, 1.5, 'CH001', -5.0, 7.0, 4, 2);

-- Insert sample sections
INSERT INTO sections (section_number, current_capacity, maximum_capacity, minimum_capacity, current_temperature, minimum_temperature, product_type_id, warehouse_id) VALUES
(1, 50, 100, 20, 5.0, -2.0, 1, 1),
(2, 30, 80, 15, 4.0, -3.0, 2, 2),
(3, 20, 50, 10, 7.0, -5.0, 3, 1),
(4, 40, 100, 20, 6.0, -4.0, 4, 2);

-- Insert sample employees
INSERT INTO employees (id_card_number, first_name, last_name, warehouse_id) VALUES
('E001', 'Alice', 'Johnson', 1),
('E002', 'Bob', 'Smith', 2);

-- Insert sample buyers
INSERT INTO buyers (id_card_number, first_name, last_name) VALUES
('B001', 'Charlie', 'Brown'),
('B002', 'Diana', 'White');

-- Insert sample product batches
INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES
('BATCH001', 500, 5.0, '2025-01-15 12:00:00', 1000, '2025-01-10', 8, 3.0, 1, 1),
('BATCH002', 300, 4.0, '2025-01-18 12:00:00', 800, '2025-01-12', 9, 2.0, 2, 2);

-- Insert sample product records
INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES
('2025-01-05 12:00:00', 2.50, 3.00, 1),
('2025-01-06 12:00:00', 1.80, 2.30, 2),
('2025-01-07 12:00:00', 0.90, 1.50, 3),
('2025-01-08 12:00:00', 5.00, 6.50, 4);

-- Insert sample inbound orders
INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES
('2025-01-05 12:00:00', 'IN001', 1, 1, 1),
('2025-01-06 12:00:00', 'IN002', 2, 2, 2);

-- Insert sample purchase orders
INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, carrier_id, order_status_id, warehouse_id) VALUES
('PO001', '2025-01-05 12:00:00', 'TRK001', 1, 1, 1, 1),
('PO002', '2025-01-06 12:00:00', 'TRK002', 2, 2, 2, 2);

-- Insert sample product records for tracking prices
INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES
('2025-01-05 12:00:00', 2.50, 3.00, 1),
('2025-01-06 12:00:00', 1.80, 2.30, 2),
('2025-01-07 12:00:00', 0.90, 1.50, 3),
('2025-01-08 12:00:00', 5.00, 6.50, 4);