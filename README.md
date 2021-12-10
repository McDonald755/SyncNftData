# SyncNftData
##DDL

DROP TABLE IF EXISTS NFT_DATA;
CREATE TABLE `NFT_DATA` (
`ID` int(11) NOT NULL AUTO_INCREMENT,
`created_time` datetime DEFAULT '2000-01-01 00:00:00',
`updated_time` datetime DEFAULT '2000-01-01 00:00:00',
`token_id` longtext,
`token_uri` longtext,
`owner` varchar(64) DEFAULT NULL,
`oracle_add` varchar(64) DEFAULT NULL,
`token_approval` longtext,
PRIMARY KEY (`ID`),
KEY `owner` (`owner`),
KEY `oracle` (`oracle_add`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

DROP TABLE IF EXISTS ORACLE_DATA;
CREATE TABLE `ORACLE_DATA` (
`ID` int(11) NOT NULL AUTO_INCREMENT,
`created_time` datetime DEFAULT NULL,
`updated_time` datetime DEFAULT NULL,
`address` varchar(64) DEFAULT NULL,
`token_symbol` varchar(255) DEFAULT NULL,
`token_name` varchar(255) DEFAULT NULL,
`approval_all` longtext,
PRIMARY KEY (`ID`)
) ENGINE=InnoDB  CHARSET=utf8mb4
