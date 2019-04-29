/*
 * @Author: rayou
 * @Date: 2019-04-15 20:45:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-18 23:22:22
 */
package output

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestDfs(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	dfsPipe := &OutputFastDfs{}

	data := NewCollectFastDfsData()
	data.Add("test2", "test.txt", []byte("test info"))
	dfsPipe.Reset("http://172.17.87.202:8080", "test")
	dfsPipe.CollectData(data.ToCollectData())
	// fmt.Printf("%+v\n", data)
	t.Log("done")

}
