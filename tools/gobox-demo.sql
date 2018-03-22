-- Adminer 4.2.5 MySQL dump

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
  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `edit_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `name` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `status` (`status`),
  KEY `add_time` (`add_time`),
  KEY `edit_time` (`edit_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

INSERT INTO `demo` (`id`, `add_time`, `edit_time`, `name`, `status`) VALUES
(1,	'2018-03-20 08:24:15',	'2018-03-20 08:24:15',	'aa',	0),
(2,	'2017-10-31 11:52:29',	'2018-02-01 09:56:09',	'tbvgtt',	3),
(10,	'2017-10-26 17:06:30',	'2017-10-26 17:06:30',	'a',	0),
(11,	'2017-10-26 17:06:30',	'2017-10-26 17:06:30',	'b',	1),
(12,	'2017-11-02 10:45:28',	'2017-11-02 10:50:19',	'abcdefg',	3),
(24,	'2017-11-09 17:09:51',	'2017-11-09 17:09:51',	'ttttttttt',	1),
(25,	'2017-11-09 17:10:09',	'2017-11-09 17:10:09',	'tttt',	1),
(27,	'2017-11-09 17:12:36',	'2017-11-09 17:12:36',	'ttttt',	1),
(32,	'2018-02-01 09:55:50',	'2018-02-01 09:55:50',	'ttttaaaa',	1),
(461,	'2017-10-27 08:20:49',	'2017-10-27 08:20:49',	'test',	1);

DROP TABLE IF EXISTS `id_gen`;
CREATE TABLE `id_gen` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `max_id` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

INSERT INTO `id_gen` (`id`, `name`, `max_id`) VALUES
(1,	'demo',	60);

-- 2018-03-22 01:56:21
