CREATE TABLE `requests` (
  `url` varchar(16) NOT NULL,
  `time` bigint(16) NOT NULL,
  `method` varchar(8) NOT NULL DEFAULT '',
  `ip` varchar(48) NOT NULL DEFAULT '',
  `post` text NOT NULL,
  `headers` text NOT NULL,
  `body` text NOT NULL,
  `referer` varchar(128) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `urls` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `url` varchar(16) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `url` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
