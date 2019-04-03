/*
 * @Author: rayou
 * @Date: 2019-04-03 18:02:56
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 18:19:19
 */

package pipeline

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

type pipelineManager struct {
	pipeline_map map[string]Pipeline
	sync.Mutex
}

var Manager pipelineManager = pipelineManager{
	pipeline_map: make(map[string]Pipeline),
}

func (self *pipelineManager) Register(pipe Pipeline) error {
	self.Lock()
	defer self.Unlock()
	name := pipe.GetName()
	elem := pipelien_map[name]
	if elem != nil {
		return errors.New("Repeat Pipeline" + name)
	}
	self.pipeline_map[name] = pipe
	return nil
}

func (self *pipelineManager) GetPipeline(name string) Pipeline {
	return self.pipeline_map[name]
}

// 初始化所有管道
func (self *pipelineManager) InitPipeline() error {
	for _, pipe := range self.pipeline_map {
		ret := pipe.Init()
		if ret != nil {
			log.Error("init pipe erros,msg:" + ret.Error())
			return ret
		}
	}
}

// 停止所有管道
func 