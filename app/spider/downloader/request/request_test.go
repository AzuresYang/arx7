/*
 * @Author: rayou
 * @Date: 2019-03-26 21:44:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-22 23:46:28
 */
package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/AzuresYang/arx7/config"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func buildCfg() *config.CrawlerTask {
	cfg := &config.CrawlerTask{
		RedisAddr:     "193.112.68.221:6379",
		RedisPassword: "Redis@2019416",
	}
	return cfg
}

func getConn() redis.Conn {
	pwd_option := redis.DialPassword("Redis@2019416")
	c, err := redis.Dial("tcp", "193.112.68.221:6379", pwd_option)
	if err != nil {
		fmt.Printf("error, dial redis fail:%s\n", err.Error())
		return nil
	}
	return c
}

func TestReqSerilize(t *testing.T) {
	// var src_req = NewArxRequest("http://www.baidu.com")
	// ret_json := src_req.Serialize()
	// fmt.Println("ret is:", ret_json)

	t.Log("gg lllllll")
	// t.Log(ret_json)
}

func TestReqManagerPushQueue(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	req_mgr := RequestMgr
	if req_mgr == nil {
		t.Error("req Mgr is nil")
		return
	}
	cfg := buildCfg()
	fmt.Printf("ready init\n")
	req_mgr.Init(cfg)
	fmt.Printf("ready add req")
	var src_req = NewArxRequest("http://www.xbiquge.la/paihangbang/")
	req_mgr.AddNeedGrabRequest(src_req)
	return
	new_req := req_mgr.GetRequest(2 * time.Second)
	if new_req == nil {
		t.Error("no req")
		return
	}
	fmt.Printf("ret is:%+v", new_req)

	t.Log("gogogo")
}

func TestAddGet(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	req_mgr := RequestMgr
	if req_mgr == nil {
		t.Error("req Mgr is nil")
		return
	}
	cfg := buildCfg()
	req_mgr.Init(cfg)
	var src_req = NewArxRequest("middle")
	req_mgr.AddNeedGrabRequest(src_req)
	src_req.Url = "hige"
	src_req.Priority = ARXREQ_PRIORITY_HIGH
	req_mgr.AddNeedGrabRequest(src_req)
	src_req.Url = "low"
	src_req.Priority = ARXREQ_PRIORITY_LOW
	req_mgr.AddNeedGrabRequest(src_req)
	new_req := req_mgr.GetRequest(2 * time.Second)
	if new_req == nil {
		t.Error("no req")
		return
	}
	fmt.Printf("ret is:%+v", new_req)

}
func TestDelKey(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	c := getConn()
	defer c.Close()
	c.Do("del", "0::ArxReqUnique")
	t.Log("gogogo")
}

func TestRedis(t *testing.T) {
	pwd_option := redis.DialPassword("Redis@2019416")
	c, err := redis.Dial("tcp", "193.112.68.221:6379", pwd_option)
	defer c.Close()
	if err != nil {
		t.Errorf("redis dial error:%s", err.Error())
	}
	unique_set := "set_test"
	unique_id := "2399"
	value := unique_id
	redis_url_key := "list_test"
	id_exit, err := redis.Bool(c.Do("SISMEMBER", unique_set, unique_id))
	if err != nil {
		t.Errorf("get memeber fail:%s", err.Error())
	}
	fmt.Printf("id_exit:%+v\n", id_exit)
	if !id_exit {
		c.Send("SADD", unique_set, unique_id)
		c.Send("RPUSH", redis_url_key, value)
		c.Flush()
		reply, _ := c.Receive()
		fmt.Printf("reply1:%+v\n", reply)
		reply, _ = c.Receive()
		if reply == nil {
			t.Error("push error")
		}
		fmt.Printf("reply1:%+v\n", reply)

	}
	t.Log("succ")
}

func TestGengerateData(t *testing.T) {
	c := getConn()
	defer c.Close()
	list := []string{
		"list_1",
		"list_2",
		"list_3",
		"list_4",
	}
	for _, l := range list {
		for i := 0; i < 2; i++ {
			value := fmt.Sprintf("%s::%d", l, i)
			reply, _ := c.Do("RPUSH", l, value)
			fmt.Printf("[%s]%d, %+v\n", l, i, reply)
		}
	}
}
func TestGetRequest(t *testing.T) {
	c := getConn()
	defer c.Close()
	list := []string{
		"list_1",
		"list_2",
		"list_3",
		"list_4",
	}
	for i := 0; i < 6; i++ {
		reply, _ := redis.Values(c.Do("BLPOP", list[0], list[1], list[2], list[3], 1))
		// fmt.Printf("reply:%+v, err:%+v\n", reply, err)
		if reply == nil {
			fmt.Printf("no data\n")
			return
		}
		for _, v := range reply {
			fmt.Printf("reply %+v\n", string(v.([]byte)))
		}
	}

	t.Log("succ")
}
