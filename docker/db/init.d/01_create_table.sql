CREATE TABLE `user` (
  `user_id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ユーザー名',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'メールアドレス',
  `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'パスワード',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_UN` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
