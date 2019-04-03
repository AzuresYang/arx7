package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Printf("%#v\n", CrawlerCfg)
	err := WriteToFile(CrawlerCfg)
	if err != nil {
		t.Error("config wirte fail,errmsg" + err.Error())
	} else {
		t.Log("write succ")
	}

}
