/*
 * @Author: rayou
 * @Date: 2019-04-11 19:59:07
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 23:04:10
 */

package config

const (
	default_userName        = "root"
	default_password        = "mysql5722"
	default_ip              = "127.0.0.1"
	default_port            = "3306"
	default_dbName          = "monitor_info"
	default_charset         = "utf8"
	default_maxConnLifetime = 100 // 数据库最大连接时间
	default_maxIdleConns    = 10  // 数据库最大闲置连接数
	default_taskNum         = 3   // 负责保存的线程数目
)

type MysqlConfig struct {
	Ip              string
	Port            string
	UserName        string
	Password        string
	Charset         string
	DbName          string
	MaxConnLifetime int
	MaxIdleConns    int
	TaskNum         int
}

func NewMysqlConfig() *MysqlConfig {
	cfg := &MysqlConfig{
		Ip:              default_ip,
		Port:            default_port,
		UserName:        default_userName,
		Password:        default_password,
		Charset:         default_charset,
		DbName:          default_dbName,
		MaxConnLifetime: default_maxConnLifetime,
		MaxIdleConns:    default_maxIdleConns,
		TaskNum:         default_taskNum,
	}
	return cfg
}
