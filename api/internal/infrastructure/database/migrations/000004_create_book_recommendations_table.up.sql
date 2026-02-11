CREATE TABLE IF NOT EXISTS book_recommendations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    recommendations_json TEXT NOT NULL,
    created_at DATETIME(6) NULL,
    expires_at DATETIME(6) NULL,
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
