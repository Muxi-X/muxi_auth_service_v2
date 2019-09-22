CREATE DATABASE `muxisite_auth`;

USE `muxisite_auth`;

CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `default` tinyint(1) DEFAULT NULL,
  `permissions` int(11) DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `ix_roles_default` (`default`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

INSERT INTO `roles` (`name`, `default`, `permissions`)
VALUES ('Moderator', 0, 14);

INSERT INTO `roles` (`name`, `default`, `permissions`)
VALUES ('Administrator', 0, 255);

INSERT INTO `roles` (`name`, `default`, `permissions`)
VALUES ('User', 1, 6);

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(164) DEFAULT NULL,
  `info` text,
  `username` varchar(164) DEFAULT NULL,
  `avatar_url` text,
  `personal_blog` text,
  `github` text,
  `flickr` text,
  `weibo` text,
  `zhihu` text,
  `password_hash` varchar(164) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  `birthday` varchar(164) DEFAULT NULL,
  `group` varchar(164) DEFAULT NULL,
  `hometown` varchar(164) DEFAULT NULL,
  `left` tinyint(1) DEFAULT NULL,
  `timejoin` varchar(164) DEFAULT NULL,
  `timeleft` varchar(164) DEFAULT NULL,
  `reset_t` varchar(164) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=254 DEFAULT CHARSET=utf8;