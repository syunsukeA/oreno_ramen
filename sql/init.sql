DROP TABLE IF EXISTS `users`;
create table `users` (
    `user_id`         BIGINT(20) AUTO_INCREMENT,
    `user_name`       VARCHAR(36) NOT NULL,
    `password`        VARCHAR(36) NOT NULL,
    `created_at`      datetime  default current_timestamp,
    -- `updated_at` timestamp default current_timestamp on update current_timestamp
    -- UNIQUE KEY uq_keys(user_name),
    PRIMARY KEY (`user_id`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

INSERT INTO users (user_name, password) VALUES ('test_user', 'pass');