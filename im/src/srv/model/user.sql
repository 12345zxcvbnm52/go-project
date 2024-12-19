CREATE TABLE `users` (
         `id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL ,
         `avatar` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
         `username` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
         `mobile` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
         `password` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
         `status` tinyint COLLATE utf8mb4_unicode_ci DEFAULT 1,
         `gender` varchar(6) COLLATE utf8mb4_unicode_ci DEFAULT '男',
         `created_at` timestamp NULL DEFAULT NULL,
         `updated_at` timestamp NULL DEFAULT NULL,
         PRIMARY KEY (`id`),
         constraint idx_uniq_name unique(`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;