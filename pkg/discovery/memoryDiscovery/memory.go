package memoryDiscovery

import (
	"context"
	"errors"
	"movie-app.com/pkg/discovery"
	"sync"
	"time"
)

/* This implementation can be used in tests or simple applications running on a single server. The
implementation is based on a combination of a map data structure and sync.RWMutex, allowing
reads and writes to the map concurrently. In the map, we store serviceInstance structures
containing the instance address and the last time of a successful health check for it, which can be set
by calling a ReportHealthyState function. In the ServiceAddresses function, we only
return instances with successful health checks from within the last 5 seconds. */

type serviceName string

// instanceID represents a unique identifier for a service instance.
type instanceID string

// Registry defines an in-memory service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

// serviceInstance represents an individual service instance.
type serviceInstance struct {
	hostPort   string    // The host and port of the service instance.
	lastActive time.Time // The last active timestamp of the instance.
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: make(map[serviceName]map[instanceID]*serviceInstance),
	}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceIDSTR string, serviceNameSTR string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	serviceNameValue := serviceName(serviceNameSTR)
	instanceIDValue := instanceID(instanceIDSTR)

	// Initialize service map if it does not exist
	if _, ok := r.serviceAddrs[serviceNameValue]; !ok {
		r.serviceAddrs[serviceNameValue] = make(map[instanceID]*serviceInstance)
	}

	// Add or update the service instance
	r.serviceAddrs[serviceNameValue][instanceIDValue] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}

	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceIDSTR string, serviceNameSTR string) error {
	r.Lock()
	defer r.Unlock()
	serviceNameValue := serviceName(serviceNameSTR)
	instanceIDValue := instanceID(instanceIDSTR)
	// If the service does not exist, return early
	if _, ok := r.serviceAddrs[serviceNameValue]; !ok {
		return nil
	}

	// Remove the specific service instance
	delete(r.serviceAddrs[serviceNameValue], instanceIDValue)

	// If no instances remain for the service, remove the service entry
	if len(r.serviceAddrs[serviceNameValue]) == 0 {
		delete(r.serviceAddrs, serviceNameValue)
	}

	return nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceIDSTR string, serviceNameSTR string) error {

	r.Lock()
	defer r.Unlock()

	serviceNameValue := serviceName(serviceNameSTR)
	instanceIDValue := instanceID(instanceIDSTR)

	if _, ok := r.serviceAddrs[serviceNameValue]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceNameValue][instanceIDValue]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceNameValue][instanceIDValue].lastActive =
		time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceNameSTR string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	serviceNameValue := serviceName(serviceNameSTR)

	if len(r.serviceAddrs[serviceNameValue]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[serviceNameValue] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
