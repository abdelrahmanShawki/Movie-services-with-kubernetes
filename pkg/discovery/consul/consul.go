package consul

import (
	"context"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"movie-app.com/pkg/discovery"
	"strconv"
	"strings"
)

type registry struct {
	client *consul.Client
}

func NewRegistry(addr string) (*registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &registry{client: client}, nil
}

// Register creates a service record in the registry.
func (r *registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of<host>:<port>, example: localhost:8081")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check:   &consul.AgentServiceCheck{CheckID: instanceID, TTL: "5s"},
	})
}

// Deregister removes a service record from the
// registry.
func (r *registry) Deregister(ctx context.Context, instanceID string, _ string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service.
func (r *registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry.
func (r *registry) ReportHealthyState(instanceID string, _ string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
