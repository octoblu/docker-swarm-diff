package reality

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/engine-api/types"
)

// Instance represents an instance of a service running
// on a server
type Instance struct {
	container types.Container
}

// NewInstance constructs a new instance given a container
func NewInstance(container types.Container) *Instance {
	return &Instance{container: container}
}

func (instance *Instance) String() string {
	state := instance.container.State

	serviceName := instance.container.Labels["com.docker.swarm.service.name"]

	taskName := instance.container.Labels["com.docker.swarm.task.name"]
	parts := strings.Split(taskName, ".")
	slot := parts[1]
	if _, err := strconv.Atoi(slot); err != nil {
		slot = "0"
	}

	return fmt.Sprintf("%v: %v.%v", state, serviceName, slot)
}
