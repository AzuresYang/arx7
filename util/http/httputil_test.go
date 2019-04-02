/*
 * @Author: rayou
 * @Date: 2019-03-27 19:02:59
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-27 19:07:46
 */
package http

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpHeader(t *testing.T) {
	header := make(http.Header)
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	header.Set("Accept-Encoding", "gzip,deflate,sdch")
	header_str := GetHeaderString(&header)
	fmt.Println("headers:", header_str)
	t.Log(header_str)

}
