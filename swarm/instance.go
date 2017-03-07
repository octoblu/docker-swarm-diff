package swarm

import (
	"fmt"

	"github.com/docker/docker/api/types/swarm"
)

// Instance is a service instance on a serveer
type Instance struct {
	serviceName string
	task        swarm.Task
}

// NewInstance returns a new instance
func NewInstance(serviceName string, task swarm.Task) *Instance {
	return &Instance{serviceName: serviceName, task: task}
}

// Key returns a key that can be used to uniquelly identify this instance
func (instance *Instance) Key() string {
	return fmt.Sprintf("%v.%v", instance.serviceName, instance.task.Slot)
}

// ServerID returns the ID of the server this
// instance is running on
func (instance *Instance) ServerID() string {
	return instance.task.NodeID
}

func (instance *Instance) String() string {
	return fmt.Sprintf("%v: %v.%v:", instance.serviceName, instance.task.Slot, instance.task.Status.State)
}
