drop table if exists `monitor_data`;
create table if not exists `monitor_data`(
	`id` bigint primary key not null auto_increment,
	`svcid` int NOT NULL,
    `metric` int not NULL,
    `classfy` int not NULL default 0,
    `value` bigint not NULL,
    `ip` varchar(20),
    `time` bigint,
    index(`svcid`, `metric`,`time`, `classfy`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;