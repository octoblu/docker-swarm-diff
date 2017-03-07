package server

// Server defines a common interface
// between reality and swarm servers
type Server interface {
	// ServiceInstances should return service instances
	// that are presumed to be running on this server
	ServiceInstances() ([]ServiceInstance, error)

	// String returns a string representation of the Server
	String() string
}

// ServiceInstance is an instance of a service presumed to
// be running on a server
type ServiceInstance interface {
	// String returns a string representation of the instance
	String() string
}
