package processor

import (
	"testing"
)

func TestManager(t *testing.T) {
	procer := NewDefaultProcessor()
	Manager.Register(&procer)
	temp_procer := Manager.GetProcessor(procer.GetName())
	// temp_procer.Free()
	if temp_procer != nil {
		t.Log("get procer succ")
	}
	t.Log(procer.GetName())
}
