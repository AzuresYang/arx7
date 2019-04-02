/*
 * @Author: rayou
 * @Date: 2019-04-02 22:25:27
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 22:31:08
 */

package processor

type (
	ProcessorManager interface {
		Register(Processor) error
		GetProcessor(string) Processor
	}
	processor_manager struct {
		processor_map map[string]Processor
	}
)

var manager processor_manager{
	processor_map: make(map[string]Processor)
}

func(self *processor_manager) Register(name string) error{
	
}
