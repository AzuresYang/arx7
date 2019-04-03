/*
 * @Author: rayou
 * @Date: 2019-04-02 22:25:27
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 23:42:44
 */

package processor

import (
	"errors"
	"sync"
)

type (
	ProcessorManager interface {
		Register(Processor) error
		GetProcessor(string) Processor
	}
	processor_manager struct {
		processor_map map[string]Processor
		sync.Mutex
	}
)

var Manager processor_manager = processor_manager{
	processor_map: make(map[string]Processor),
}

func (self *processor_manager) Register(processor Processor) error {
	self.Lock()
	defer self.Unlock()
	name := processor.GetName()
	elem := self.processor_map[name]
	if elem != nil {
		return errors.New("Repeat Processor:" + name)
	}
	self.processor_map[name] = processor
	return nil
}

func (self *processor_manager) GetProcessor(name string) Processor {
	procer := self.processor_map[name]
	if procer == nil {
		return nil
	}
	return procer.GetOneProcessor()
}
