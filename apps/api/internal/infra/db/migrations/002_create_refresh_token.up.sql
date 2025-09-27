CREATE TABLE refresh_tokens (
    jti         CHAR(26) NOT NULL PRIMARY KEY,
    user_id     CHAR(26) NOT NULL,
    revoked_at  DATETIME NULL,
    expires_at  DATETIME NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_user_active (user_id, revoked_at, expires_at),
    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;