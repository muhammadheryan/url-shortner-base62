-- migrate:up
CREATE TABLE url (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    short_url VARCHAR(20) DEFAULT "",
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);


-- migrate:down
DROP TABLE url;
