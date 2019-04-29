package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/AzuresYang/arx7/app/arxmonitor/monitorHandler"
	"github.com/AzuresYang/arx7/app/status"
	"github.com/AzuresYang/arx7/config"
	"github.com/AzuresYang/arx7/util/stringUtil"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

const (
	default_redis_read_timeout   time.Duration = 2 * time.Second
	default_redis_write_timeout  time.Duration = 2 * time.Second
	default_redis_urlqueue_name  string        = "ArxUrlQueue"
	default_redis_url_unique_set string        = "ArxReqUnique"
	default_max_queue_len        int           = 100 // 管道中最多存在未处理请求
)

type RequestManager struct {
	wait_push_req_queue chan *ArxRequest //等待推送到主服务的req无锁队列
	redisPool           *redis.Pool
	isDistribute        bool // 是否是分布式，也就Redis
	taskId              uint32
}

var RequestMgr = &RequestManager{
	isDistribute: true,
}

func NewRequestManager() *RequestManager {
	return &RequestManager{
		isDistribute: true,
	}
}
func (self *RequestManager) Init(cfg *config.CrawlerTask) error {
	code_info := "RequestManager.Init"
	self.taskId = cfg.TaskId
	self.wait_push_req_queue = make(chan *ArxRequest, default_max_queue_len)
	// 单机模式不需要链接redis
	if !self.isDistribute {
		return nil
	}
	err := self.testRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Errorf("[%s]test redis conn fail:%s", code_info, err.Error())
		return err
	}
	pwd_option := redis.DialPassword(cfg.RedisPassword)
	write_option := redis.DialWriteTimeout(default_redis_write_timeout)
	read_option := redis.DialReadTimeout(default_redis_read_timeout)
	self.redisPool = &redis.Pool{
		MaxIdle: 5, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", cfg.RedisAddr, pwd_option, write_option, read_option)
		},
	}
	log.Infof("[%s]redis init succ, addr:%s", code_info, cfg.RedisAddr)
	return nil
}

// 测试redis联通性
func (self *RequestManager) testRedis(redisAddr string, pwd string) error {
	pwd_option := redis.DialPassword(pwd)
	c, err := redis.Dial("tcp", redisAddr, pwd_option)
	if c != nil {
		c.Close()
		return nil
	}
	return err
}

// 清空redis队列
func (self *RequestManager) ClearRedis(redisAddr string, pwd string) error {
	unique_set := self.getRedisUniqueSetName()
	pwd_option := redis.DialPassword(pwd)
	c, err := redis.Dial("tcp", redisAddr, pwd_option)
	if c != nil {
		_, derr := c.Do("del", unique_set)
		if derr != nil {
			return errors.New(fmt.Sprintf("del unique set fail:%s", derr.Error()))
		}
		c.Close()
	}
	return err
}

func (self *RequestManager) SetKeyValue(redisAddr string, pwd string, key string, value string) error {
	pwd_option := redis.DialPassword(pwd)
	c, err := redis.Dial("tcp", redisAddr, pwd_option)
	if c != nil {
		_, derr := c.Do("set", key, value)
		if derr != nil {
			return errors.New(fmt.Sprintf("set redis fails:%s", derr.Error()))
		}
		c.Close()
	}
	return err
}

func (self *RequestManager) Start() {

}

func (self *RequestManager) Stop() {

}

// 获取一个请求， 等可以连上redis之后，从redis中获取
func (self *RequestManager) getRequestWhenOneInstance(timeout time.Duration) (req *ArxRequest) {
	code_info := "RequestManager.getRequestWhenOneInstance"
	click := time.After(timeout)
	log.Debugf("[%s] get request", code_info)
	select {
	case req = <-self.wait_push_req_queue:
		return req
	case <-click:
		return nil
	}
}
func (self *RequestManager) getRedisUrlQueueName(p ArxReqPriority) string {
	priority := "MIDDLE"
	switch p {
	case ARXREQ_PRIORITY_HIGH:
		priority = "HIGH"
	case ARXREQ_PRIORITY_LOW:
		priority = "LOW"
	case ARXREQ_PRIORITY_MIDDLE:
		priority = "MIDDLE"
	}
	// taskID::urlqueue::priority
	return fmt.Sprintf("%d::%s::%s", self.taskId, default_redis_urlqueue_name, priority)
}

func (self *RequestManager) getRedisUniqueSetName() string {
	return fmt.Sprintf("%d::%s", self.taskId, default_redis_url_unique_set)
}

func (self *RequestManager) GetRequest(timeout time.Duration) *ArxRequest {
	code_info := reflect.TypeOf(self).String()
	if !self.isDistribute {
		return self.getRequestWhenOneInstance(timeout)
	}
	// 分布式模式下， 从redis获取链接
	conn := self.redisPool.Get()
	defer conn.Close()
	// 超时设置，未满1s的，设置成1s(redis 超时单位)
	timeout_sec := timeout.Seconds()
	if timeout_sec < 1 {
		timeout_sec = 1
	}
	high_queue := self.getRedisUrlQueueName(ARXREQ_PRIORITY_HIGH)
	middle_queue := self.getRedisUrlQueueName(ARXREQ_PRIORITY_MIDDLE)
	lower_queue := self.getRedisUrlQueueName(ARXREQ_PRIORITY_LOW)
	// 从高到低的优先级开始获取链接
	reply, err := redis.Values(conn.Do("BLPOP", high_queue, middle_queue, lower_queue, timeout_sec))
	if err != nil {
		log.Errorf("[%s] blpop request for redis fail:%s", code_info, err.Error())
		return nil
	}
	list := string(reply[0].([]byte))
	log.Debugf("[%s] get request from:%s", code_info, list)
	req_bytes := reply[1].([]byte)
	req := new(ArxRequest)
	err = json.Unmarshal(req_bytes, req)
	if err != nil {
		log.Errorf("[%s]unserialize request from redis fail:%s", code_info, err.Error())
		return nil
	}
	log.Debugf("[%s] get request from redis url:%s", code_info, req.Url)

	monitorHandler.AddOne(status.MONI_SYS_REQUEST_GET)
	return req
}

// 计算请求的唯一标识符
func (self *RequestManager) calRequestUniqueId(req *ArxRequest) string {
	// 唯一标识符计算方法： url_methed_postdata_tempStrMap
	s := fmt.Sprintf("%s|%s|%s|%+v", req.Url, req.Method, req.PostData, req.TempStrMap)
	hash_code := stringUtil.Hash(s)
	return fmt.Sprintf("%d", hash_code)
}

// 无锁队列入列， 需要维护一个去重的URL队列
// 已经存在的，直接添加存在
func (self *RequestManager) AddNeedGrabRequest(req *ArxRequest) error {
	code_info := "RequestManager.AddRequest"
	log.Debugf("[%s]Add new url:%s", code_info, req.Url)
	// 单机模式不需要链接redis
	if !self.isDistribute {
		self.wait_push_req_queue <- req
		return nil
	}
	// unique_id := self.calRequestUniqueId(req)
	c := self.redisPool.Get()
	defer c.Close()

	value, err := req.Serialize()
	if err != nil {
		log.Errorf("[%s] req serialize fail.", code_info)
		return err
	}
	unique_set := self.getRedisUniqueSetName()
	unique_id := self.calRequestUniqueId(req)

	// 查看请求之前是否已经存在。已存在则不添加请求到redis里
	is_exit, _ := redis.Bool(c.Do("SISMEMBER", unique_set, unique_id))
	if is_exit {
		log.Debugf("[%s]request has exit:%s", code_info, req.Url)
		return nil
	}
	// 添加一个req到redis，需要有两个操作
	// 添加req的唯一标识符到redis中
	// 添加req到redis中，
	redis_url_key := self.getRedisUrlQueueName(req.Priority)
	c.Send("RPUSH", redis_url_key, value)
	c.Send("SADD", unique_set, unique_id)
	c.Flush()
	reply1, _ := c.Receive()
	reply2, _ := c.Receive()
	if reply1 == nil || reply2 == nil {
		err_msg := fmt.Sprintf("add request to redis error.Push unique_id:%+v, Add request:%+v", reply1, reply2)
		log.Errorf("[%s]%s", code_info, err_msg)
		return errors.New(err_msg)
	}
	log.Debugf("[%s]Add new url succ:%s", code_info, req.Url)
	monitorHandler.AddOne(status.MONI_SYS_REQUEST_GET)
	return nil
}

// 下载失败的链接处理
func (self *RequestManager) AddDownLoadFailReqeust(req *ArxRequest, msg string) bool {
	return true
}

// 记录下载成功链接
func (self *RequestManager) AddDownloadSuccReq(req *ArxRequest) {

}

func (self *RequestManager) AddProcessReq(req *ArxRequest) {

}

// 处理成功链接
