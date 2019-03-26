/*
 * @Author: rayou
 * @Date: 2019-03-25 22:23:56
 * @Last Modified by:   rayou
 * @Last Modified time: 2019-03-25 22:23:56
 */

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	fileName := "test_log.txt"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Error("open file error.")
	}
	log.SetOutput(logFile)
	log.SetLevel(log.InfoLevel)
	log.Error("hello", "ERROR")
	log.Info("info")
	log.Debug("debug")
}
