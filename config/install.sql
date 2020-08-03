-- MySQL dump 10.13  Distrib 8.0.19, for osx10.14 (x86_64)
--
-- Host: localhost    Database: nging
-- ------------------------------------------------------
-- Server version	8.0.19

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
-- Table structure for table `nging_access_log`
--

DROP TABLE IF EXISTS `nging_access_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_access_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `vhost_id` int unsigned NOT NULL COMMENT '虚拟主机ID',
  `remote_addr` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'IP地址',
  `x_real_ip` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'IP',
  `x_forward_for` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'IP',
  `local_addr` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '本机地址',
  `elapsed` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '耗时(秒)',
  `host` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Header中的Host, 一般会是域名',
  `user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户',
  `time_local` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '本地时间格式',
  `minute` varchar(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'H:i',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'GET POST HEAD OPTIONS PUT DELETE...',
  `uri` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网址',
  `version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'HTTP/1.0, HTTP/1.1 ...',
  `status_code` mediumint unsigned NOT NULL DEFAULT '0' COMMENT '状态码',
  `body_bytes` bigint unsigned NOT NULL DEFAULT '0' COMMENT '响应body字节数',
  `referer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源网址',
  `user_agent` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户代理',
  `hit_status` int unsigned NOT NULL DEFAULT '0' COMMENT '缓存服务器命中状态',
  `scheme` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'https/http',
  `brower_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器名',
  `brower_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器类型(spider/mobile/pc)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

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
  `topic` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '通知专题',
  `recipient_id` int unsigned NOT NULL DEFAULT '0' COMMENT '收信账号',
  `disabled` enum('Y','N') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `recipient_id` (`recipient_id`),
  KEY `topic` (`topic`,`disabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='报警收信专题关联';
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
  UNIQUE KEY `code` (`code`)
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='验证码';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_collector_export`
--

DROP TABLE IF EXISTS `nging_collector_export`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_collector_export` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `page_root` int unsigned NOT NULL DEFAULT '0' COMMENT '根页面ID',
  `page_id` int unsigned NOT NULL DEFAULT '0' COMMENT '页面ID',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '组ID',
  `mapping` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字段映射',
  `dest` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '目标',
  `dest_type` enum('API','DSN','dbAccountID') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'dbAccountID' COMMENT '目标类型',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '方案名',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `exported` int unsigned DEFAULT '0' COMMENT '最近导出时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'N' COMMENT '是否禁用',
  PRIMARY KEY (`id`),
  KEY `idx_page_id` (`page_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='导出规则';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_collector_export_log`
--

DROP TABLE IF EXISTS `nging_collector_export_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_collector_export_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `page_id` int unsigned NOT NULL DEFAULT '0' COMMENT '页面规则ID',
  `export_id` int unsigned NOT NULL DEFAULT '0' COMMENT '导出方案ID',
  `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '结果',
  `status` enum('idle','start','success','failure') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'idle' COMMENT '状态',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `export_id` (`export_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='导出日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_collector_group`
--

DROP TABLE IF EXISTS `nging_collector_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_collector_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名',
  `type` enum('page','export') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'page' COMMENT '类型(page-页面规则组;export-导出规则组)',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='采集规则组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_collector_history`
--

DROP TABLE IF EXISTS `nging_collector_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_collector_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父ID',
  `page_id` int unsigned NOT NULL DEFAULT '0' COMMENT '页面ID',
  `page_parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父页面ID',
  `page_root_id` int unsigned NOT NULL DEFAULT '0' COMMENT '入口页面ID',
  `has_child` enum('N','Y') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否有子级',
  `url` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面网址',
  `url_md5` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面网址MD5',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面标题',
  `content` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面内容MD5',
  `rule_md5` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规则标识MD5',
  `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '采集到的数据',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `exported` int unsigned NOT NULL DEFAULT '0' COMMENT '最近导出时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `url_md5` (`url_md5`),
  KEY `idx_page_id` (`page_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='采集历史';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_collector_page`
--

DROP TABLE IF EXISTS `nging_collector_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_collector_page` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父级规则',
  `root_id` int unsigned NOT NULL DEFAULT '0' COMMENT '根页面ID',
  `has_child` enum('N','Y') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否有子级',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '规则组',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '规则名',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `enter_url` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '入口网址模板(网址一行一个)',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `browser` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器',
  `type` enum('list','content') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'content' COMMENT '页面类型',
  `scope_rule` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面区域规则',
  `duplicate_rule` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'none' COMMENT '去重规则(url-判断网址;rule-判断规则是否改过;content-判断网页内容是否改过;none-不去重)',
  `content_type` enum('html','json','text') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'html' COMMENT '内容类型',
  `charset` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '字符集',
  `timeout` int unsigned NOT NULL DEFAULT '0' COMMENT '超时时间(秒)',
  `waits` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '等待时间范围(秒),例如2-8',
  `proxy` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '代理地址',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `group_id` (`group_id`),
  KEY `parent_id` (`parent_id`),
  KEY `root_id` (`root_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='采集页面';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_collector_rule`
--

DROP TABLE IF EXISTS `nging_collector_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_collector_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `page_id` int unsigned NOT NULL COMMENT '页面ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '保存匹配结果的变量名',
  `rule` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '规则',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'string' COMMENT '数据类型',
  `filter` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '过滤器',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`id`),
  KEY `page_id` (`page_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='页面中的元素采集规则';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_command`
--

DROP TABLE IF EXISTS `nging_command`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_command` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `description` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '说明',
  `command` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '命令',
  `work_directory` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作目录',
  `env` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '环境变量(一行一个，格式为：var1=val1)',
  `created` int unsigned NOT NULL COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `remote` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否执行远程SSH命令',
  `ssh_account_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'SSH账号ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='指令集';
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
  `value` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '值',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '简介',
  `type` enum('id','text','url','html','image','video','audio','file','json','list') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'text' COMMENT '值类型(list-以半角逗号分隔的值列表)',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `encrypted` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否加密',
  PRIMARY KEY (`key`,`group`),
  KEY `g` (`group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_db_account`
--

DROP TABLE IF EXISTS `nging_db_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_db_account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
  `uid` int unsigned NOT NULL COMMENT 'UID',
  `engine` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'mysql' COMMENT '数据库引擎',
  `host` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'localhost:3306' COMMENT '服务器地址',
  `user` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'root' COMMENT '用户名',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '数据库名称',
  `options` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '其它选项(JSON)',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='数据库账号';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_db_sync`
--

DROP TABLE IF EXISTS `nging_db_sync`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_db_sync` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(120) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '方案名',
  `source_account_id` int unsigned NOT NULL DEFAULT '0' COMMENT '源数据库账号ID',
  `dsn_source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '同步源',
  `destination_account_id` int unsigned NOT NULL DEFAULT '0' COMMENT '目标数据库账号ID',
  `dsn_destination` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '目标数据库',
  `tables` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '要同步的表',
  `skip_tables` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '要跳过的表',
  `alter_ignore` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '要忽略的列、索引、外键',
  `drop` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '删除待同步数据库中多余的字段、索引、外键 ',
  `mail_to` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '发送邮件',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='数据表同步方案';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_db_sync_log`
--

DROP TABLE IF EXISTS `nging_db_sync_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_db_sync_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `sync_id` int unsigned NOT NULL COMMENT '同步方案ID',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '结果报表',
  `change_tables` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '被更改的表',
  `change_table_num` int unsigned NOT NULL DEFAULT '0' COMMENT '被更改的表的数量',
  `elapsed` bigint unsigned NOT NULL DEFAULT '0' COMMENT '总共耗时',
  `failed` int unsigned NOT NULL DEFAULT '0' COMMENT '失败次数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='数据表同步日志';
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
  `project` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '项目名称',
  `table_id` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '关联表数据id',
  `table_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '关联表名称',
  `field_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '关联表字段名',
  `sort` bigint NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(1-已审核/0-未审核)',
  `category_id` int unsigned NOT NULL DEFAULT '0' COMMENT '分类ID',
  `used_times` int unsigned NOT NULL DEFAULT '0' COMMENT '被使用的次数',
  PRIMARY KEY (`id`),
  KEY `category_id` (`category_id`),
  KEY `owner_id_owner_type` (`owner_id`,`owner_type`),
  KEY `table_id_table_name` (`table_id`,`table_name`)
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
  `table_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '表名称',
  `field_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字段名',
  `file_ids` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件id列表',
  `embedded` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否(Y/N)为内嵌文件',
  PRIMARY KEY (`id`),
  UNIQUE KEY `tid_field_name_table_name` (`table_id`,`field_name`,`table_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='嵌入文件';
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
  UNIQUE KEY `file_id_size_flag` (`file_id`,`size`),
  UNIQUE KEY `save_path` (`save_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='图片文件缩略图';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_forever_process`
--

DROP TABLE IF EXISTS `nging_forever_process`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_forever_process` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '添加人ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `command` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '命令',
  `workdir` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作目录',
  `env` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '环境变量',
  `args` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '命令参数',
  `pidfile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'PID记录文件',
  `logfile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '日志记录文件',
  `errfile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '错误记录文件',
  `respawn` int unsigned NOT NULL DEFAULT '1' COMMENT '重试次数(进程被外部程序结束后自动启动)',
  `delay` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '延迟启动(例如1ms/1s/1m/1h)',
  `ping` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '心跳时间(例如1ms/1s/1m/1h)',
  `pid` int NOT NULL DEFAULT '0' COMMENT 'PID',
  `status` enum('started','running','stopped','restarted','exited','killed','idle') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'idle' COMMENT '进程运行状态',
  `debug` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT 'DEBUG',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `error` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '错误信息',
  `lastrun` int unsigned NOT NULL DEFAULT '0' COMMENT '上次运行时间',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `enable_notify` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否启用通知',
  `notify_email` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知人列表',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='持久进程';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_frp_client`
--

DROP TABLE IF EXISTS `nging_frp_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_frp_client` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `server_addr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0.0.0.0',
  `server_port` int unsigned NOT NULL DEFAULT '7000',
  `http_proxy` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `pool_count` int unsigned NOT NULL DEFAULT '1',
  `tcp_mux` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '必须跟服务端一致',
  `user` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `dns_server` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `login_fail_exit` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y',
  `protocol` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'tcp' COMMENT 'tcp or kcp',
  `heartbeat_interval` bigint NOT NULL DEFAULT '30',
  `heartbeat_timeout` bigint NOT NULL DEFAULT '90',
  `log_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'console',
  `log_way` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'console' COMMENT 'console or file',
  `log_level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'info',
  `log_max_days` int unsigned NOT NULL DEFAULT '3',
  `token` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `admin_addr` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `admin_port` int unsigned NOT NULL DEFAULT '0',
  `admin_user` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `admin_pwd` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `start` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '要启动的代理节点名称，留空代表全部，多个用半角逗号隔开',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `uid` int unsigned NOT NULL DEFAULT '0',
  `group_id` int unsigned NOT NULL DEFAULT '0',
  `type` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'web',
  `created` int unsigned NOT NULL DEFAULT '0',
  `updated` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='FRP客户端设置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_frp_group`
--

DROP TABLE IF EXISTS `nging_frp_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_frp_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='FRP服务组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_frp_server`
--

DROP TABLE IF EXISTS `nging_frp_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_frp_server` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `tcp_mux` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y',
  `addr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0.0.0.0',
  `port` int unsigned NOT NULL DEFAULT '7000',
  `udp_port` int unsigned NOT NULL DEFAULT '0',
  `kcp_port` int unsigned NOT NULL DEFAULT '0',
  `proxy_addr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0.0.0.0',
  `vhost_http_port` int unsigned NOT NULL DEFAULT '0',
  `vhost_http_timeout` bigint unsigned NOT NULL DEFAULT '0',
  `vhost_https_port` int unsigned NOT NULL DEFAULT '0',
  `log_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'console',
  `log_way` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'console' COMMENT 'console or file',
  `log_level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'info',
  `log_max_days` int unsigned NOT NULL DEFAULT '3',
  `token` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `auth_timeout` bigint unsigned NOT NULL DEFAULT '900',
  `subdomain_host` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `max_ports_per_client` bigint NOT NULL DEFAULT '0',
  `max_pool_count` int unsigned NOT NULL DEFAULT '5',
  `heart_beat_timeout` int unsigned NOT NULL DEFAULT '90',
  `user_conn_timeout` int unsigned NOT NULL DEFAULT '10',
  `dashboard_addr` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0.0.0.0',
  `dashboard_port` int unsigned NOT NULL DEFAULT '0',
  `dashboard_user` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'admin',
  `dashboard_pwd` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'admin',
  `allow_ports` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `uid` int unsigned NOT NULL DEFAULT '0',
  `group_id` int unsigned NOT NULL DEFAULT '0',
  `created` int unsigned NOT NULL DEFAULT '0',
  `updated` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='FRP服务器设置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_ftp_user`
--

DROP TABLE IF EXISTS `nging_ftp_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_ftp_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `password` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `banned` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁止连接',
  `directory` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '授权目录(一行一个) ',
  `ip_whitelist` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'IP白名单(一行一个) ',
  `ip_blacklist` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'IP黑名单(一行一个) ',
  `created` int unsigned NOT NULL COMMENT '创建时间 ',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户组',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='FTP用户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_ftp_user_group`
--

DROP TABLE IF EXISTS `nging_ftp_user_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_ftp_user_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名称',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `banned` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁止组内用户连接',
  `directory` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '授权目录',
  `ip_whitelist` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'IP白名单(一行一个)',
  `ip_blacklist` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'IP黑名单(一行一个)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='FTP用户组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_kv`
--

DROP TABLE IF EXISTS `nging_kv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_kv` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '关键字',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '元素值',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型标识',
  `sort` int NOT NULL COMMENT '排序',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `child_key_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '子键类型(number/text...)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key_type` (`key`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='键值数据';
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
-- Table structure for table `nging_ssh_user`
--

DROP TABLE IF EXISTS `nging_ssh_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_ssh_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL COMMENT 'UID',
  `host` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'localhost' COMMENT '主机名',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口',
  `charset` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '字符集',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'root' COMMENT '用户名',
  `password` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '账号名称',
  `options` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '其它选项(JSON)',
  `private_key` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '私钥内容',
  `passphrase` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '私钥口令',
  `protocol` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '连接协议',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '说明',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '组',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='数据库账号';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_ssh_user_group`
--

DROP TABLE IF EXISTS `nging_ssh_user_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_ssh_user_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='SSH账号组';
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
  KEY `idx_uid` (`uid`),
  KEY `idx_group_id` (`group_id`)
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
  KEY `idx_uid` (`uid`)
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
  KEY `idx_task_id` (`task_id`,`created`)
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
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '盐值',
  `safe_pwd` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '安全密码',
  `avatar` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `gender` enum('male','female','secret') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'secret' COMMENT '性别(male-男;female-女;secret-保密)',
  `last_login` int unsigned NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '状态',
  `online` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否在线',
  `role_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID(多个用“,”分隔开)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `file_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上传文件总大小',
  `file_num` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上传文件数量',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
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
  `perm_cmd` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '指令集权限(多个用“,”分隔开)',
  `perm_action` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '行为权限(多个用“,”分隔开)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_user_u2f`
--

DROP TABLE IF EXISTS `nging_user_u2f`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_user_u2f` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL COMMENT '用户ID',
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '签名',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '扩展设置',
  `created` int unsigned NOT NULL COMMENT '绑定时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_type` (`uid`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='两步验证';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_vhost`
--

DROP TABLE IF EXISTS `nging_vhost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_vhost` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网站名称',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '组',
  `domain` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '域名',
  `root` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '网站物理路径',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `setting` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '设置',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否停用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='虚拟主机';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nging_vhost_group`
--

DROP TABLE IF EXISTS `nging_vhost_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nging_vhost_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='虚拟主机组';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-08-03 13:49:45
