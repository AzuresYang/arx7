package monitor

import (
	"testing"
	// "github.com/AzuresYang/arx7/app/pipeline"

	log "github.com/sirupsen/logrus"
	// "github.com/AzuresYang/arx7/config"
	// "github.com/AzuresYang/arx7/util/record"
)

func TestMonitor(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	InitMonitorHandler("127.0.0.1", 8001, 5542)
	for i := 0; i < 50; i++ {
		AddOne(uint32(i))
	}
	for {
		if IfStop() {
			return
		}
	}
}
