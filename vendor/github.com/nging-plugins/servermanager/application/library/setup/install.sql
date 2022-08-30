-- MySQL dump 10.13  Distrib 8.0.30, for macos12.4 (x86_64)
--
-- Host: 127.0.0.1    Database: nging
-- ------------------------------------------------------
-- Server version	8.0.30

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
  UNIQUE KEY `forever_process_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='持久进程';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-08-27 13:20:02
