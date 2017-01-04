-- Adminer 4.2.5 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `ftp_user`;
CREATE TABLE `ftp_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(120) NOT NULL COMMENT '用户名',
  `password` varchar(150) NOT NULL COMMENT '密码',
  `banned` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁止连接',
  `directory` varchar(500) NOT NULL COMMENT '授权目录(一行一个) ',
  `ip_whitelist` text NOT NULL COMMENT 'IP白名单(一行一个) ',
  `ip_blacklist` text NOT NULL COMMENT 'IP黑名单(一行一个) ',
  `created` int(10) unsigned NOT NULL COMMENT '创建时间 ',
  `updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `group_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户组',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='FTP用户';


DROP TABLE IF EXISTS `ftp_user_group`;
CREATE TABLE `ftp_user_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '组名称',
  `created` int(10) unsigned NOT NULL COMMENT '创建时间',
  `updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `disabled` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `banned` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁止组内用户连接',
  `directory` varchar(500) NOT NULL COMMENT '授权目录',
  `ip_whitelist` text NOT NULL COMMENT 'IP白名单(一行一个)',
  `ip_blacklist` text NOT NULL COMMENT 'IP黑名单(一行一个)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='FTP用户组';


DROP TABLE IF EXISTS `vhost`;
CREATE TABLE `vhost` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `domain` text NOT NULL COMMENT '域名',
  `root` varchar(500) NOT NULL COMMENT '网站物理路径',
  `created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `setting` text NOT NULL COMMENT '设置',
  `disabled` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否停用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='虚拟主机';


-- 2016-12-11 11:16:40
