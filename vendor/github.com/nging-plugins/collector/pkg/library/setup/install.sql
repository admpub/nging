-- MySQL dump 10.13  Distrib 8.0.26, for macos11.3 (x86_64)
--
-- Host: 127.0.0.1    Database: nging
-- ------------------------------------------------------
-- Server version	8.0.26

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
  KEY `collector_export_page_id` (`page_id`)
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
  KEY `collector_export_log_export_id` (`export_id`)
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
  KEY `collector_group_uid` (`uid`)
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
  UNIQUE KEY `collector_history_url_md5` (`url_md5`),
  KEY `collector_history_page_id` (`page_id`)
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
  KEY `collector_page_uid` (`uid`),
  KEY `collector_page_group_id` (`group_id`),
  KEY `collector_page_parent_id` (`parent_id`),
  KEY `collector_page_root_id` (`root_id`)
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
  KEY `collector_rule_page_id` (`page_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='页面中的元素采集规则';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-12-15 11:40:08
