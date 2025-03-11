CREATE TABLE IF NOT EXISTS users (
    id INt AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    id_rol INT NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    state BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP NULL,
    failed_attempts INT DEFAULT 0,
    last_failed_attempt TIMESTAMP NULL,
    password_changed_at TIMESTAMP NULL,
    force_password_change BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS roles (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  is_root BOOLEAN DEFAULT FALSE
  created_by INT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (created_by) REFERENCES users(id)
);

ALTER TABLE users
ADD FOREIGN KEY (id_rol) REFERENCES roles(id);

CREATE TABLE IF NOT EXISTS modules (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(255) UNIQUE NOT NULL,
  menu BOOLEAN DEFAULT FALSE,
  parent_id INT NULL,
  order_index INT NOT NULL,
  FOREIGN KEY (parent_id) REFERENCES Module(id)
);

CREATE TABLE IF NOT EXISTS rol_modules (
  id INT AUTO_INCREMENT PRIMARY KEY,
  role_id INT NOT NULL,
  module_id INT NOT NULL,
  can_view BOOLEAN DEFAULT FALSE,
  can_write BOOLEAN DEFAULT FALSE,
  can_edit BOOLEAN DEFAULT FALSE,
  can_delete BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by INT NULL,
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (module_id) REFERENCES modules(id),
  FOREIGN KEY (created_by) REFERENCES users(id),
  UNIQUE KEY unique_role_module (role_id, module_id)
);

-- Tabla para intentos de acceso
CREATE TABLE IF NOT EXISTS LoginAttempt (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN DEFAULT FALSE,
    user_agent TEXT,
    failure_reason VARCHAR(255)
);

-- Tabla para bloqueos temporales
CREATE TABLE IF NOT EXISTS AccountLock (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    locked_until TIMESTAMP NOT NULL,
    reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Tabla para sesiones
CREATE TABLE IF NOT EXISTS AuthSession (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_user INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (id_user) REFERENCES users(id)
);

-- Media (imágenes, archivos)
CREATE TABLE IF NOT EXISTS Media (
  id INT AUTO_INCREMENT PRIMARY KEY,
  file_name VARCHAR(255) NOT NULL,
  file_path VARCHAR(255) NOT NULL,
  mime_type VARCHAR(20) NOT NULL,
  size DECIMAL(10, 6) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);
-- Categorías con estructura jerárquica
CREATE TABLE IF NOT EXISTS Categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  slug VARCHAR(220) UNIQUE NOT NULL,
  description TEXT,
  parent_id INT NULL,  -- Para subcategorías
  id_media INT NULL,   -- Imagen principal
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (parent_id) REFERENCES Categories(id),
  FOREIGN KEY (id_media) REFERENCES Media(id)
);
-- Tags
CREATE TABLE IF NOT EXISTS Tags (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  code VARCHAR(100) UNIQUE NOT NULL,
  active BOOLEAN, 
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL, 
);

--Taxes
CREATE TABLE IF NOT EXISTS TaxCategory (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  status BOOLEAN, 
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL, 
);

--Products
CREATE TABLE IF NOT EXISTS Products (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  slug VARCHAR(220) UNIQUE NOT NULL,
  description TEXT,
  status BOOLEAN,
  id_category INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL, 
  FOREIGN KEY (id_category) REFERENCES Categories(id)
);
--N:M ProductosTags
CREATE TABLE IF NOT EXISTS N_Products_Tags (
  id_tag INT,
  id_product INT,
  PRIMARY KEY (id_tag, id_product),
  FOREIGN KEY (id_tag) REFERENCES Tags(id)
  FOREIGN KEY (id_product) REFERENCES Products(id)
);
--ProductAtributte
CREATE TABLE IF NOT EXISTS ProductAttributes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  type VARCHAR(100) UNIQUE NOT NULL,
  description TEXT, 
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL, 
);

--N:M ProductAttribute Products
CREATE TABLE IF NOT EXISTS N_ProductAttirbute_Products (
  id_product INT,
  id_product_attribute INT,
  PRIMARY KEY (id_product, id_product_attribute),
  FOREIGN KEY (id_product) REFERENCES Products(id)
  FOREIGN KEY (id_product_attribute) REFERENCES ProductAttributes(id)
);
--AtributteValue
CREATE TABLE IF NOT EXISTS AtributteValue (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  id_attribute_product INT,
  FOREIGN KEY (id_attribute_product) REFERENCES ProductAttributes(id)
);

--ProductVariant
CREATE TABLE IF NOT EXISTS ProductVariant (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  sku VARCHAR(100) UNIQUE NOT NULL,
  stock INT,
  id_product INT,
  id_tax_category INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL, 
  FOREIGN KEY (id_product) REFERENCES Products(id)
  FOREIGN KEY (id_tax_category) REFERENCES TaxCategory(id)
);

--N:M ProductoVariant AtributeValue
CREATE TABLE IF NOT EXISTS N_ProductVariant_AttributeValue (
  id_attribute_value INT,
  id_product_variant INT,
  PRIMARY KEY (id_attribute_value, id_product_variant),
  FOREIGN KEY (id_attribute_value) REFERENCES AtributteValue(id)
  FOREIGN KEY (id_product_variant) REFERENCES ProductVariant(id)
);
--N:M ProductoVariant Media
CREATE TABLE IF NOT EXISTS N_ProductVariant_Media (
  id_product_variant INT,
  id_media INT,
  PRIMARY KEY (id_media, id_product_variant),
  FOREIGN KEY (id_media) REFERENCES Media(id)
  FOREIGN KEY (id_product_variant) REFERENCES ProductVariant(id)
);

-- Zonas
CREATE TABLE IF NOT EXISTS Zones (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- TaxZones
CREATE TABLE IF NOT EXISTS TaxZones (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_zone INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_zone) REFERENCES Zones(id)
);

-- TaxRates
CREATE TABLE IF NOT EXISTS TaxRates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_tax_category INT NOT NULL,
    id_tax_zone INT NOT NULL,
    percentage DECIMAL(5,2) NOT NULL,
    is_default BOOLEAN DEFAULT false,
    status BOOLEAN DEFAULT true,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_tax_category) REFERENCES TaxCategory(id),
    FOREIGN KEY (id_tax_zone) REFERENCES TaxZones(id)
);

-- TaxRules
CREATE TABLE IF NOT EXISTS TaxRules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_tax_rate INT NOT NULL,
    priority INT DEFAULT 0,
    status BOOLEAN DEFAULT true,
    min_amount DECIMAL(10,2),
    max_amount DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_tax_rate) REFERENCES TaxRates(id)
);

-- PriceLists
CREATE TABLE IF NOT EXISTS PriceLists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT false,
    priority INT NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Prices
CREATE TABLE IF NOT EXISTS Prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_price_list INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    starts_at TIMESTAMP NOT NULL,
    ends_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_price_list) REFERENCES PriceLists(id)
);
-- N:M ProductoVariante Precio
CREATE TABLE IF NOT EXISTS N_ProductVariant_Prices (
    id_product_variant INT,
    id_price INT,
    is_active BOOLEAN,
    PRIMARY KEY (id_price, id_product_variant),
    FOREIGN KEY (id_price) REFERENCES Prices(id)
    FOREIGN KEY (id_product_variant) REFERENCES ProductVariant(id)
);

-- Currencies
CREATE TABLE IF NOT EXISTS Currencies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(5) NOT NULL UNIQUE,
    symbol VARCHAR(5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- ActualCurrency
CREATE TABLE IF NOT EXISTS ActualCurrency (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    is_base BOOLEAN DEFAULT false,
    ratio DECIMAL(10,6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

--N:M ProductoVariant Media
CREATE TABLE IF NOT EXISTS N_ActualCurrency_Currency (
  id_currency INT,
  id_actual_currency INT,
  PRIMARY KEY (id_currency, id_actual_currency),
  FOREIGN KEY (id_currency) REFERENCES Currencies(id)
  FOREIGN KEY (id_actual_currency) REFERENCES ActualCurrency(id)
)

--Country
CREATE TABLE IF NOT EXISTS Countries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    flag VARCHAR(100) NOT NULL,
    translations VARCHAR(100) NOT NULL,
    iso2 VARCHAR(5) NOT NULL,
    phone_code VARCHAR(10) NOT NULL,
    id_currency INT,
    FOREIGN KEY (id_currency) REFERENCES Currencies(id)
);
--City
CREATE TABLE IF NOT EXISTS Cities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    id_country INT,
    FOREIGN KEY (id_country) REFERENCES Countries(id)
);

-- N:M City Zones
CREATE TABLE IF NOT EXISTS N_Cities_Zones (
    id_city INT,
    id_zone INT,
    PRIMARY KEY (id_city, id_zone),
    FOREIGN KEY (id_city) REFERENCES Cities(id)
    FOREIGN KEY (id_zone) REFERENCES Zones(id)
);


-- Customers
CREATE TABLE IF NOT EXISTS Customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    active BOOLEAN DEFAULT true,
    is_guest BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
-- Customers
CREATE TABLE IF NOT EXISTS CustomersGroup (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_price INT,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_price) REFERENCES Prices(id),
);

-- N:M cUSTOMERS CUSTOMERSgROUP
CREATE TABLE IF NOT EXISTS N_Customer_CustomerGroup (
    id_customer INT,
    id_customer_group INT,
    PRIMARY KEY (id_customer, id_customer_group),
    FOREIGN KEY (id_customer) REFERENCES Customers(id)
    FOREIGN KEY (id_customer_group) REFERENCES CustomersGroup(id)
);


-- ShippingAddresses
CREATE TABLE IF NOT EXISTS ShippingAddresses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_customer INT NOT NULL,
    id_city INT NOT NULL,
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    postal_code VARCHAR(10),
    province VARCHAR(100),
    district VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_customer) REFERENCES Customers(id),
    FOREIGN KEY (id_city) REFERENCES Cities(id)
);

-- ShippingMethods
CREATE TABLE IF NOT EXISTS ShippingMethods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    code VARCHAR(50) NOT NULL,
    calculator_type VARCHAR(50),
    price_default DECIMAL(10,2),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- ShippingRules
CREATE TABLE IF NOT EXISTS ShippingRules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_shipping_method INT NOT NULL,
    max_value DECIMAL(10,2),
    min_value DECIMAL(10,2),
    price DECIMAL(10,2) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_shipping_method) REFERENCES ShippingMethods(id)
);


-- Carts
CREATE TABLE IF NOT EXISTS Carts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_customer INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    session_id VARCHAR(255),
    total_amount DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_customer) REFERENCES Customers(id)
);

-- CartItems
CREATE TABLE IF NOT EXISTS CartItems (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_cart INT NOT NULL,
    id_product_variant INT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_cart) REFERENCES Carts(id),
    FOREIGN KEY (id_product_variant) REFERENCES ProductVariant(id)
);


-- EligibilityCheckers
CREATE TABLE IF NOT EXISTS EligibilityCheckers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_shipping_method INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    rule JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_shipping_method) REFERENCES ShippingMethods(id)
);

-- N:M Zones-ShippingMethods
CREATE TABLE IF NOT EXISTS N_Zones_ShippingMethod (
    id_shipping_method INT,
    id_zone INT,
    PRIMARY KEY (id_shipping_method, id_zone),
    FOREIGN KEY (id_shipping_method) REFERENCES ShippingMethods(id),
    FOREIGN KEY (id_zone) REFERENCES Zones(id)
);

-- PaymentMethods
CREATE TABLE IF NOT EXISTS PaymentMethods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Orders
CREATE TABLE IF NOT EXISTS Orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_customer INT NOT NULL,
    id_shipping_address INT NOT NULL,
    id_shipping_method INT NOT NULL,
    id_currency INT NOT NULL,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    shipping_cost DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    total_items INT NOT NULL,
    tracking_number VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (id_customer) REFERENCES Customers(id),
    FOREIGN KEY (id_shipping_address) REFERENCES ShippingAddresses(id),
    FOREIGN KEY (id_shipping_method) REFERENCES ShippingMethods(id),
    FOREIGN KEY (id_currency) REFERENCES Currencies(id)
);

-- OrderItems
CREATE TABLE IF NOT EXISTS OrderItems (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_order INT NOT NULL,
    id_product_variant INT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_order) REFERENCES Orders(id),
    FOREIGN KEY (id_product_variant) REFERENCES ProductVariant(id)
);

-- Payments
CREATE TABLE IF NOT EXISTS Payments (
    id INT AUTO_INCREMENT,
    id_order INT NOT NULL,
    id_payment_method INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id, id_order),
    FOREIGN KEY (id_order) REFERENCES Orders(id),
    FOREIGN KEY (id_payment_method) REFERENCES PaymentMethods(id)
);

-- ProductDimensions (para cálculos de envío)
CREATE TABLE IF NOT EXISTS ProductDimensions (
    id_product_variant INT PRIMARY KEY,
    length DECIMAL(10,2),
    width DECIMAL(10,2),
    height DECIMAL(10,2),
    weight DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_product_variant) REFERENCES ProductVariant(id)
);

--INDICES
CREATE INDEX idx_product_slug ON Products(slug);
CREATE INDEX idx_order_number ON Orders(order_number);
CREATE INDEX idx_customer_email ON Customers(email);