package processor

import (
	"fmt"
	"regexp"
	"testing"
)

func TestManager(t *testing.T) {
	procer := NewDefaultProcessor()
	Manager.Register(&procer)
	temp_procer := Manager.GetProcessor(procer.GetName())
	// temp_procer.Free()
	if temp_procer != nil {
		t.Log("get procer succ")
	}
	t.Log(procer.GetName())
}

func TestGetXiaoshuo(t *testing.T) {
	line := "<dd><a href='/17/17377/8705057.html' >书友天图分香纵论贼道三痴历史文</a></dd>"
	// <a href="http://www.xbiquge.la/([\d]+/[\d]+)/">(.+)</a></li>
	flysnowRegexp := regexp.MustCompile(`<dd><a href='(.+)' >(.+)</a></dd>`)
	params := flysnowRegexp.FindStringSubmatch(line)
	for i, it := range params {
		fmt.Printf("[%d]%s\n", i, it)
	}

}
