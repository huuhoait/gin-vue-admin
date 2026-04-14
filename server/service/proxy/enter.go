package proxy

import "time"

// ServiceGroup holds singleton proxy clients for downstream services.
type ServiceGroup struct {
	Core  *CoreProxy
	Order *OrderProxy
}

// NewServiceGroup initialises proxy clients from the given config values.
func NewServiceGroup(coreURL, orderURL string, timeoutSec int) *ServiceGroup {
	timeout := time.Duration(timeoutSec) * time.Second
	return &ServiceGroup{
		Core:  NewCoreProxy(coreURL, timeout),
		Order: NewOrderProxy(orderURL, timeout),
	}
}
