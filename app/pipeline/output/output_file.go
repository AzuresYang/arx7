/*
 * @Author: rayou
 * @Date: 2019-04-05 22:05:31
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-05 22:35:38
 */

package output

import (
	"github.com/AzuresYang/arx7/app/pipeline"
	log "github.com/sirupsen/logrus"
)

type OutputFile struct {
}

func (self *OutputFile) Init() error {
	return nil
}

func (self *OutputFile) CollectData(data *pipeline.CollectData) {
	log.Debug("collect data")
}

func (self *OutputFile) Stop() {

}

func (self *OutputFile) GetName() string {
	return "OutputFile"
}
