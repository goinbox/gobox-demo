-- Adminer 4.6.2 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP DATABASE IF EXISTS `gobox-demo`;
CREATE DATABASE `gobox-demo` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */;
USE `gobox-demo`;

DROP TABLE IF EXISTS `demo`;
CREATE TABLE `demo` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `add_time` timestamp NOT NULL DEFAULT current_timestamp(),
  `edit_time` timestamp NOT NULL DEFAULT current_timestamp(),
  `name` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `status` tinyint(4) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `status` (`status`),
  KEY `add_time` (`add_time`),
  KEY `edit_time` (`edit_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

INSERT INTO `demo` (`id`, `add_time`, `edit_time`, `name`, `status`) VALUES
(63,	'2019-12-09 07:28:34',	'2019-12-09 07:32:52',	'bbb',	1),
(65,	'2019-12-09 07:32:52',	'2019-12-09 07:32:52',	'aaa',	1);

DROP TABLE IF EXISTS `id_gen`;
CREATE TABLE `id_gen` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `max_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

INSERT INTO `id_gen` (`id`, `name`, `max_id`) VALUES
(1,	'demo',	65);

-- 2020-01-21 02:48:41
