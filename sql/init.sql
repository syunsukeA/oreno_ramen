USE `oreno_ramen_db`;
create table `users` (
    `user_id`         BIGINT(20) AUTO_INCREMENT,
    `profile_img`     text,  -- ToDo: 適切なデータ型決める
    `username`        VARCHAR(36) NOT NULL,
    `password`        VARCHAR(36) NOT NULL,
    `created_at`      datetime  default current_timestamp,
    `updated_at`      timestamp default current_timestamp on update current_timestamp,
    UNIQUE KEY (`username`),
    PRIMARY KEY (`user_id`)
);

create table `shops` (
    `shop_id`           VARCHAR(36) NOT NULL UNIQUE, -- HotPepperのPK
    `user_id`           BIGINT(20) NOT NULL,
    `shopname`          VARCHAR(36) NOT NULL,
    `bookmark`          TINYINT(1) NOT NULL, -- 0or1の値 (ToDo: Goの方でBooleanの数値的扱いを確認するべきかも)
    `created_at`        datetime  default current_timestamp,
    `updated_at`        timestamp default current_timestamp on update current_timestamp,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`)
);

create table `reviews` (
    `review_id`         BIGINT(20) AUTO_INCREMENT,
    `user_id`           BIGINT(20) NOT NULL,
    `shop_id`           VARCHAR(36) NOT NULL, -- HotPepperのPK
    `shopname`          VARCHAR(36) NOT NULL,
    `content`           text, -- ToDo: データ型の選定
    `evaluate`          INT CHECK (evaluate >= 0 AND evaluate <=5),
    `review_img`        text,  -- ToDo: 適切なデータ型決める, ToDo: 画像が複数設定可能な場合どうする？
    `created_at`        datetime  default current_timestamp,
    `updated_at`        timestamp default current_timestamp on update current_timestamp,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`),
    PRIMARY KEY (`review_id`)
);

-- Initial data for users table
INSERT INTO users (username, password) VALUES ('test_user', 'pass');
INSERT INTO users (username, password) VALUES ('syunsuke', 'hoge');
INSERT INTO users (username, password) VALUES ('guest1', '0120');

-- Initial data for reviews table
INSERT INTO reviews (user_id, shop_id, shopname) VALUES (1, '000000', 'test_review');