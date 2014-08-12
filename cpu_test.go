package i6502

import (
	"testing"
)

func TestNewCpu(t *testing.T) {
	cpu, err := NewCpu()

	if err != nil {
		t.Errorf("Expected NewCPU() to not raise an error")
	}

	if cpu == nil {
		t.Errorf("Expected NewCPU() to create a new CPU instance")
	}
}
