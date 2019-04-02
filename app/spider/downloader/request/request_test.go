/*
 * @Author: rayou
 * @Date: 2019-03-26 21:44:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-26 22:16:56
 */
package request

import (
	"fmt"
	"testing"
	"time"
)

func TestReqSerilize(t *testing.T) {
	var src_req = NewArxRequest("http://www.baidu.com")
	ret_json := src_req.Serialize()
	// fmt.Println("ret is:", ret_json)

	t.Log("gg lllllll")
	t.Log(ret_json)
}

func TestReqManagerPushQueue(t *testing.T) {
	req_mgr := RequestManager{}
	req_mgr.Init(2)
	var src_req = NewArxRequest("http://www.baidu.com")
	req_mgr.AddNeedGrabRequest(&src_req, 2*time.Second)
	new_req := req_mgr.GetRequest()
	if new_req == nil {
		t.Error("req is lost")
		return
	}
	ret_json := new_req.Serialize()
	fmt.Println("ret is:", ret_json)
	t.Log("gogogo")
}
