-- PayTheFly crypto payment orders table
CREATE TABLE IF NOT EXISTS paythefly_orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    serial_no VARCHAR(128) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    quota INT NOT NULL,
    amount VARCHAR(64) NOT NULL,
    deadline BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    tx_hash VARCHAR(128) DEFAULT NULL,
    wallet VARCHAR(128) DEFAULT NULL,
    paid_value VARCHAR(64) DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_serial_no (serial_no),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
