drop table if exists `monitor_data`;
create table if not exists `monitor_data`(
	`svcid` int NOT NULL,
    `metric` int not NULL,
    `classfy` int not NULL default 0,
    `value` bigint not NULL,
    `ip` varchar(20),
    `time` timestamp not null,
    primary key(`svcid`, `metric`)
    
)ENGINE=InnoDB DEFAULT CHARSET=utf8;monitor_data