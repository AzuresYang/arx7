package processor

import (
	"testing"
)

func TestManager(t *testing.T) {
	procer := NewDefaultProcessor()
	Manager.Register(&procer)
	temp_procer := Manager.GetProcessor(procer.GetName())
	temp_procer.Free()
	t.Log(procer.GetName())
}
