/*
 * @Author: rayou
 * @Date: 2019-03-27 18:52:42
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-27 19:10:01
 */

package http

import (
	"encoding/json"
	"net/http"
)

func GetHeaderString(header *http.Header) string {
	json_byte, _ := json.Marshal(header)
	return string(json_byte[:])
}
