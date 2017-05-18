package swarm

import (
	"fmt"

	"github.com/docker/engine-api/types/swarm"
)

// Instance is a service instance on a serveer
type Instance struct {
	serviceName string
	service     swarm.Service
	task        swarm.Task
}

// NewInstance returns a new instance
func NewInstance(serviceName string, service swarm.Service, task swarm.Task) *Instance {
	return &Instance{serviceName: serviceName, service: service, task: task}
}

// Key returns a key that can be used to uniquelly identify this instance
func (instance *Instance) Key() string {
	if instance.task.Slot == 0 {
		return fmt.Sprintf("%v.%v", instance.serviceName, instance.task.NodeID)
	}
	return fmt.Sprintf("%v.%v", instance.serviceName, instance.task.Slot)
}

// ServerID returns the ID of the server this
// instance is running on
func (instance *Instance) ServerID() string {
	return instance.task.NodeID
}

// ShouldIgnore returns true if this instance should be ignored
// typically only used if the instance is supposed to be shutdown
// on a server
func (instance *Instance) ShouldIgnore() bool {
	if instance.service.UpdateStatus.State == swarm.UpdateStateUpdating {
		return true
	}
	return instance.task.Status.State == swarm.TaskStateShutdown
}

func (instance *Instance) String() string {
	return fmt.Sprintf("%v: %v.%v", instance.task.Status.State, instance.serviceName, instance.task.Slot)
}
