-- Create database for NFT events
CREATE DATABASE IF NOT EXISTS nft_events 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE nft_events;

-- Contract events table (will be auto-created by GORM, but provided for reference)
CREATE TABLE IF NOT EXISTS contract_events (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    event_name VARCHAR(100) NOT NULL COMMENT 'Event name',
    contract_address VARCHAR(42) NOT NULL COMMENT 'Contract address',
    tx_hash VARCHAR(66) NOT NULL UNIQUE COMMENT 'Transaction hash',
    block_number BIGINT UNSIGNED NOT NULL COMMENT 'Block number',
    block_hash VARCHAR(66) COMMENT 'Block hash',
    log_index INT UNSIGNED NOT NULL COMMENT 'Log index',
    from_address VARCHAR(42) COMMENT 'From address',
    to_address VARCHAR(42) COMMENT 'To address',
    data TEXT COMMENT 'Event data',
    topics TEXT COMMENT 'Event topics',
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Event timestamp',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation time',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update time',
    INDEX idx_event_name (event_name),
    INDEX idx_contract_address (contract_address),
    INDEX idx_block_number (block_number),
    INDEX idx_tx_hash (tx_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Contract events storage';

-- Block heights tracking table
CREATE TABLE IF NOT EXISTS block_heights (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    height BIGINT UNSIGNED NOT NULL COMMENT 'Last processed block height',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
    INDEX idx_height (height)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Block height tracking';

-- Insert initial block height record
INSERT INTO block_heights (height) VALUES (0) 
ON DUPLICATE KEY UPDATE height = VALUES(height);
