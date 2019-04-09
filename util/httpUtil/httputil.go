/*
 * @Author: rayou
 * @Date: 2019-03-27 18:52:42
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-09 20:40:29
 */

package httpUtil

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

func GetHeaderString(header *http.Header) string {
	json_byte, _ := json.Marshal(header)
	return string(json_byte[:])
}

func GetLocalIp() (ip string, e error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("counld get local ip")
}
