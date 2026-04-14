package utils

import (
	"sync"
)

// SystemEvents DefineSystemLevelEventHandle
type SystemEvents struct {
	reloadHandlers []func() error
	mu             sync.RWMutex
}

// GlobalEventmanagementDevice
var GlobalSystemEvents = &SystemEvents{}

// RegisterReloadHandler register system reload handler
func (e *SystemEvents) RegisterReloadHandler(handler func() error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.reloadHandlers = append(e.reloadHandlers, handler)
}

// TriggerReload TriggerAllRegisterofOverloadHandleFunctionNumber
func (e *SystemEvents) TriggerReload() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	for _, handler := range e.reloadHandlers {
		if err := handler(); err != nil {
			return err
		}
	}
	return nil
}
