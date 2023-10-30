-- MySQL dump 10.13  Distrib 8.1.0, for macos12.6 (x86_64)
--
-- Host: 127.0.0.1    Database: nging
-- ------------------------------------------------------
-- Server version	8.1.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `nging_alert_recipient`
--

DROP TABLE IF EXISTS `nging_alert_recipient`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_alert_recipient` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `account` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '账号',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '扩展信息(JSON)',
  `type` enum('email','webhook') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'email' COMMENT '类型',
  `platform` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '平台(dingding-钉钉;workwx-企业微信)',
  `description` varchar(500) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='报警收信人';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_alert_topic`
--

DROP TABLE IF EXISTS `nging_alert_topic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_alert_topic` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `topic` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '通知专题',
  `recipient_id` int unsigned NOT NULL DEFAULT '0' COMMENT '收信账号',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `alert_topic_recipient_id` (`recipient_id`),
  KEY `alert_topic_topic_disabled` (`topic`,`disabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='报警收信专题关联';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_cloud_backup`
--

DROP TABLE IF EXISTS `nging_cloud_backup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_cloud_backup` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '配置名',
  `source_path` varchar(200) COLLATE utf8mb4_general_ci NOT NULL COMMENT '源',
  `ignore_rule` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '忽略文件路径(正则表达式)',
  `wait_fill_completed` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否等待文件填充结束',
  `min_modify_interval` int unsigned NOT NULL DEFAULT '0' COMMENT '最小修改间隔',
  `ignore_wait_rule` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '忽略等待文件完成的规则',
  `delay` int unsigned NOT NULL DEFAULT '0' COMMENT '延后秒数',
  `storage_engine` enum('s3','sftp','ftp','webdav','smb') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 's3' COMMENT '存储引擎',
  `storage_config` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '存储引擎配置参数(JSON)',
  `dest_storage` int unsigned NOT NULL COMMENT '目标存储ID',
  `dest_path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '目标存储路径',
  `result` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '运行结果',
  `last_executed` int unsigned NOT NULL DEFAULT '0' COMMENT '最近运行时间',
  `status` enum('idle','running','failure') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'idle' COMMENT '运行状态',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `log_disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用日志',
  `log_type` enum('error','all') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'all' COMMENT '日志类型(error-仅记录报错;all-记录所有)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='云备份';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_cloud_backup_log`
--

DROP TABLE IF EXISTS `nging_cloud_backup_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_cloud_backup_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `backup_id` int unsigned NOT NULL DEFAULT '0' COMMENT '云备份规则ID',
  `backup_type` enum('full','change') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'change' COMMENT '备份方式(full-全量备份;change-文件更改时触发的备份)',
  `backup_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '需要备份本地文件',
  `remote_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '保存到远程的文件路径',
  `operation` enum('create','update','delete','none') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'none' COMMENT '操作',
  `error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '错误信息',
  `status` enum('success','failure') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'success' COMMENT '状态',
  `elapsed` int unsigned NOT NULL DEFAULT '0' COMMENT '消耗时间(毫秒)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `cloud_backup_log_backup_id_created` (`backup_id`,`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='云备份日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_cloud_storage`
--

DROP TABLE IF EXISTS `nging_cloud_storage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_cloud_storage` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `type` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'aws' COMMENT '存储类型(aws,oss,cos)',
  `key` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'key',
  `secret` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密钥(加密处理)',
  `bucket` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '存储桶',
  `endpoint` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '地域节点',
  `region` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '地区',
  `secure` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否(Y/N)HTTPS',
  `baseurl` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '资源基础网址',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='云存储账号';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_code_invitation`
--

DROP TABLE IF EXISTS `nging_code_invitation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_code_invitation` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
  `recv_uid` int unsigned NOT NULL DEFAULT '0' COMMENT '使用者',
  `code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '邀请码',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `used` int unsigned NOT NULL DEFAULT '0' COMMENT '使用时间',
  `start` int unsigned NOT NULL DEFAULT '0' COMMENT '有效时间',
  `end` int unsigned NOT NULL DEFAULT '0' COMMENT '失效时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `role_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '注册为角色(多个用“,”分隔开)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code_invitation_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='邀请码';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_code_verification`
--

DROP TABLE IF EXISTS `nging_code_verification`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_code_verification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '验证码',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所有者ID',
  `owner_type` enum('user','customer') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '所有者类型',
  `used` int unsigned NOT NULL DEFAULT '0' COMMENT '使用时间',
  `purpose` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '目的',
  `start` int unsigned NOT NULL DEFAULT '0' COMMENT '有效时间',
  `end` int unsigned NOT NULL DEFAULT '0' COMMENT '失效时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `send_method` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'mobile' COMMENT '发送方式(mobile-手机;email-邮箱)',
  `send_to` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送目标',
  PRIMARY KEY (`id`),
  KEY `code_verification_disabled_owner` (`disabled`,`owner_type`,`owner_id`,`send_method`,`send_to`,`code`,`purpose`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='验证码';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_config`
--

DROP TABLE IF EXISTS `nging_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_config` (
  `key` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '键',
  `group` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '组',
  `label` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '选项名称',
  `value` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '值',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '简介',
  `type` enum('id','text','url','html','image','video','audio','file','json','list') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'text' COMMENT '值类型(list-以半角逗号分隔的值列表)',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `encrypted` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否加密',
  PRIMARY KEY (`key`,`group`),
  KEY `config_group` (`group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_file`
--

DROP TABLE IF EXISTS `nging_file`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_file` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `owner_type` enum('user','customer') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '用户类型',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '原始文件名',
  `save_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '保存名称',
  `save_path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件保存路径',
  `view_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '查看链接',
  `ext` varchar(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件后缀',
  `mime` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件mime类型',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'image' COMMENT '文件类型',
  `size` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件大小',
  `width` int unsigned NOT NULL DEFAULT '0' COMMENT '宽度(像素)',
  `height` int unsigned NOT NULL DEFAULT '0' COMMENT '高度(像素)',
  `dpi` int unsigned NOT NULL DEFAULT '0' COMMENT '分辨率',
  `md5` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件md5',
  `storer_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件保存位置',
  `storer_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '位置ID',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '上传时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `sort` bigint NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(1-已审核/0-未审核)',
  `category_id` int unsigned NOT NULL DEFAULT '0' COMMENT '分类ID',
  `tags` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标签',
  `subdir` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '子目录',
  `used_times` int unsigned NOT NULL DEFAULT '0' COMMENT '被使用的次数',
  PRIMARY KEY (`id`),
  KEY `file_category_id` (`category_id`),
  KEY `file_view_url` (`view_url`),
  KEY `file_subdir` (`subdir`),
  KEY `file_owner_type_and_id` (`owner_type`,`owner_id`),
  KEY `file_name` (`save_name`,`name`),
  KEY `file_type` (`type`),
  KEY `file_used_times` (`used_times`),
  KEY `file_storer` (`storer_name`,`storer_id`),
  KEY `file_size` (`size` DESC),
  KEY `file_sort` (`sort`),
  KEY `file_status` (`status`),
  KEY `file_width_height` (`width` DESC,`height` DESC),
  KEY `file_updated` (`updated` DESC),
  KEY `file_created` (`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_file_embedded`
--

DROP TABLE IF EXISTS `nging_file_embedded`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_file_embedded` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `project` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '项目名',
  `table_id` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '表主键',
  `table_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '表名称',
  `field_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字段名',
  `file_ids` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件id列表',
  `embedded` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否(Y/N)为内嵌文件',
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_embedded_table_id_field_table` (`table_id`,`field_name`,`table_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='嵌入文件';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_file_moved`
--

DROP TABLE IF EXISTS `nging_file_moved`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_file_moved` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_id` bigint unsigned NOT NULL COMMENT '文件ID',
  `from` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件原路径',
  `to` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件新路径',
  `thumb_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '缩略图ID(缩略图时有效)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_moved_from` (`from`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='文件移动记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_file_thumb`
--

DROP TABLE IF EXISTS `nging_file_thumb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_file_thumb` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_id` bigint unsigned NOT NULL COMMENT '文件ID',
  `size` bigint unsigned NOT NULL COMMENT '文件大小',
  `width` int unsigned NOT NULL COMMENT '宽度(像素)',
  `height` int unsigned NOT NULL COMMENT '高度(像素)',
  `dpi` int unsigned NOT NULL DEFAULT '0' COMMENT '分辨率',
  `save_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '保存名称',
  `save_path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '保存路径',
  `view_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '访问网址',
  `used_times` int unsigned NOT NULL DEFAULT '0' COMMENT '被使用的次数',
  `md5` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '缩略图文件MD5值',
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_thumb_save_path` (`save_path`),
  UNIQUE KEY `file_thumb_file_id_size_flag` (`file_id`,`size`),
  KEY `file_thumb_view_url` (`view_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='图片文件缩略图';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_kv`
--

DROP TABLE IF EXISTS `nging_kv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_kv` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '键名',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '元素值',
  `description` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `help` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '帮助说明',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型标识',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `child_key_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'text' COMMENT '子键类型(number/text...)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `kv_key_type` (`key`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='键值数据';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_login_log`
--

DROP TABLE IF EXISTS `nging_login_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_login_log` (
  `owner_type` enum('customer','user') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '用户类型(user-后台用户;customer-前台客户)',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'session id',
  `username` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录名',
  `auth_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'password' COMMENT '认证方式',
  `errpwd` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '错误密码',
  `ip_address` varchar(46) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'ip地址',
  `ip_location` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'ip定位',
  `user_agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器代理',
  `success` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否登录成功',
  `failmsg` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '失败信息',
  `day` int unsigned NOT NULL DEFAULT '0' COMMENT '日期(Ymd)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  KEY `login_log_ip_address` (`ip_address`,`day`),
  KEY `login_log_created` (`created` DESC),
  KEY `login_log_owner` (`owner_type`,`owner_id`,`session_id`),
  KEY `login_log_success` (`success`),
  KEY `login_log_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='登录日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_sending_log`
--

DROP TABLE IF EXISTS `nging_sending_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_sending_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `sent_at` int unsigned NOT NULL COMMENT '发送时间',
  `source_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '来源ID',
  `source_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '来源类型',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `method` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'mobile' COMMENT '发送方式(mobile-手机;email-邮箱)',
  `to` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送目标',
  `provider` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送平台',
  `result` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送结果描述',
  `status` enum('success','failure','waiting','queued','none') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'waiting' COMMENT '发送状态(none-无需发送)',
  `retries` int unsigned NOT NULL DEFAULT '0' COMMENT '重试次数',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送消息内容',
  `params` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送消息参数(JSON)',
  `appointment_time` int unsigned NOT NULL DEFAULT '0' COMMENT '预约发送时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='邮件短信等发送日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_task`
--

DROP TABLE IF EXISTS `nging_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_task` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '分组ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务名称',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '任务类型',
  `description` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '时间表达式',
  `concurrent` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否支持多实例',
  `command` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '命令详情',
  `work_directory` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作目录',
  `env` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '环境变量(一行一个，格式为：var1=val1)',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `enable_notify` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否启用通知',
  `notify_email` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知人列表',
  `timeout` bigint unsigned NOT NULL DEFAULT '0' COMMENT '超时设置',
  `execute_times` int unsigned NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int unsigned NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `closed_log` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)关闭日志',
  PRIMARY KEY (`id`),
  KEY `task_uid` (`uid`),
  KEY `task_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_task_group`
--

DROP TABLE IF EXISTS `nging_task_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_task_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `cmd_prefix` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '命令前缀',
  `cmd_suffix` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '命令后缀',
  PRIMARY KEY (`id`),
  KEY `task_group_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_task_log`
--

DROP TABLE IF EXISTS `nging_task_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_task_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务输出',
  `error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '错误信息',
  `status` enum('success','timeout','failure','stop','restart') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'success' COMMENT '状态',
  `elapsed` int unsigned NOT NULL DEFAULT '0' COMMENT '消耗时间(毫秒)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `task_log_task_id_created` (`task_id`,`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_user`
--

DROP TABLE IF EXISTS `nging_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`密码',
  `salt` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`盐值',
  `safe_pwd` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`安全密码',
  `session_id` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`session id',
  `avatar` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `gender` enum('male','female','secret') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'secret' COMMENT '性别(male-男;female-女;secret-保密)',
  `last_login` int unsigned NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `login_fails` int unsigned NOT NULL DEFAULT '0' COMMENT '连续登录失败次数',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '状态',
  `online` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否在线',
  `role_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID(多个用“,”分隔开)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `file_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上传文件总大小',
  `file_num` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上传文件数量',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_user_role`
--

DROP TABLE IF EXISTS `nging_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_user_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `description` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '说明',
  `created` int unsigned NOT NULL COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_user_role_permission`
--

DROP TABLE IF EXISTS `nging_user_role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_user_role_permission` (
  `role_id` int unsigned NOT NULL COMMENT '角色ID',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '权限类型',
  `permission` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '权限值',
  UNIQUE KEY `user_role_permission_uniqid` (`role_id`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_user_u2f`
--

DROP TABLE IF EXISTS `nging_user_u2f`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_user_u2f` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL COMMENT '用户ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '签名',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '扩展设置',
  `step` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '第几步',
  `precondition` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '除了密码登录外的其它前置条件(仅step=2时有效),用半角逗号分隔',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '绑定时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_u2f_uid_type` (`uid`,`type`),
  KEY `user_u2f_step` (`step`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='两步验证';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-10-30 17:17:31
