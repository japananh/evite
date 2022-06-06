SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `users` (
   `id` int PRIMARY KEY AUTO_INCREMENT,
   `email` varchar(50) UNIQUE NOT NULL,
    `password` varchar(50) NOT NULL,
    `salt` varchar(50) NOT NULL,
    `status` smallint unsigned NOT NULL DEFAULT 1,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;
