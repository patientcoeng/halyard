package alerting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAlert struct {
	triggeredCount int
}

func (m *MockAlert) Trigger(summary string, detail string) {
	m.triggeredCount++
}

func TestManager_AddAlert(t *testing.T) {
	manager := NewManager()
	manager.AddAlert(&MockAlert{})
	assert.Equal(t, 1, len(manager.alerts))
}

func TestManager_Trigger(t *testing.T) {
	manager := NewManager()
	mockAlert := &MockAlert{}
	manager.AddAlert(mockAlert)

	manager.Trigger("test", "test")

	assert.Equal(t, 1, mockAlert.triggeredCount)
}

func TestManager_MultipleTrigger(t *testing.T) {
	manager := NewManager()
	mockAlertA := &MockAlert{}
	mockAlertB := &MockAlert{}
	manager.AddAlert(mockAlertA)
	manager.AddAlert(mockAlertB)

	manager.Trigger("test", "test")

	assert.Equal(t, 1, mockAlertA.triggeredCount)
	assert.Equal(t, 1, mockAlertB.triggeredCount)
}
