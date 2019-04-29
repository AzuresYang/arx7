/*
 * @Author: rayou
 * @Date: 2019-03-27 19:17:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-27 19:37:11
 */

package downloader

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/AzuresYang/arx7/app/spider/downloader/request"
)

func TestSimpleDownloader(t *testing.T) {
	var src_req = request.NewArxRequest("https://www.bilibili.com/read/cv2320240/")

	simple_downloader := SimpleDownloader{}

	ctx := simple_downloader.Download(nil, src_req)
	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)
	bodystr := string(body)
	fmt.Println(bodystr)
	t.Log("down")
}
