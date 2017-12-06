-- Adminer 4.2.5 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `code_invitation`;
CREATE TABLE `code_invitation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `code` varchar(40) NOT NULL COMMENT '邀请码',
  `created` int(11) unsigned NOT NULL COMMENT '创建时间',
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `recv_uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '使用者',
  `used` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '使用时间',
  `start` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '有效时间',
  `end` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '失效时间',
  `disabled` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='邀请码';


DROP TABLE IF EXISTS `code_verification`;
CREATE TABLE `code_verification` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `code` varchar(40) NOT NULL COMMENT '验证码',
  `created` int(11) unsigned NOT NULL COMMENT '创建时间',
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `used` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '使用时间',
  `purpose` varchar(40) NOT NULL COMMENT '目的',
  `start` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '有效时间',
  `end` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '失效时间',
  `disabled` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='验证码';


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


DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '分组ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(100) NOT NULL DEFAULT '' COMMENT '时间表达式',
  `concurrent` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否支持多实例',
  `command` text NOT NULL COMMENT '命令详情',
  `disabled` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `enable_notify` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否启用通知',
  `notify_email` text NOT NULL COMMENT '通知人列表',
  `timeout` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '超时设置',
  `execute_times` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `created` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务';


DROP TABLE IF EXISTS `task_group`;
CREATE TABLE `task_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(60) NOT NULL COMMENT '组名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
  `created` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务组';


DROP TABLE IF EXISTS `task_log`;
CREATE TABLE `task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` enum('success','timeout','failure','stop','restart') NOT NULL DEFAULT 'success' COMMENT '状态',
  `elapsed` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '消耗时间(毫秒)',
  `created` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务日志';


DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` char(64) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(64) NOT NULL DEFAULT '' COMMENT '盐值',
  `last_login` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` varchar(150) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `disabled` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户';


DROP TABLE IF EXISTS `user_u2f`;
CREATE TABLE `user_u2f` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) unsigned NOT NULL COMMENT '用户ID',
  `token` varchar(255) NOT NULL COMMENT '签名',
  `type` varchar(30) NOT NULL COMMENT '类型',
  `extra` text NOT NULL COMMENT '扩展设置',
  `created` int(11) unsigned NOT NULL COMMENT '绑定时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_type` (`uid`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='两步验证';


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


-- 2017-04-09 11:49:42

DROP TABLE IF EXISTS `db_account`;
CREATE TABLE `db_account` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int(11) unsigned NOT NULL COMMENT 'UID',
  `engine` varchar(30) NOT NULL DEFAULT 'mysql' COMMENT '数据库引擎',
  `host` varchar(200) NOT NULL DEFAULT 'localhost:3306' COMMENT '服务器地址',
  `user` varchar(100) NOT NULL DEFAULT 'root' COMMENT '用户名',
  `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(120) NOT NULL DEFAULT '' COMMENT '数据库名称',
  `options` text NOT NULL COMMENT '其它选项(JSON)',
  `created` int(10) unsigned NOT NULL COMMENT '创建时间',
  `updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='数据库账号';

-- 2017-10-29 12:19:00
