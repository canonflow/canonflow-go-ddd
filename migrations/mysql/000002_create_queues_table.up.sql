CREATE TABLE IF NOT EXISTS queues (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    unique_id VARCHAR(255) NOT NULL UNIQUE,
    queue VARCHAR(100) NOT NULL,
    payload JSON,
    status ENUM(
        'pending',
        'running',
        'done',
        'failed',
        'cancelled'
    ) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_queue_status (queue, STATUS),
    INDEX idx_unique_id (unique_id)
)