CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `nickname` varchar(65) NOT NULL DEFAULT '' COMMENT '昵称',
  `phone` varchar(12) NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(65) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` varchar(65) NOT NULL DEFAULT '' COMMENT '密码',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;