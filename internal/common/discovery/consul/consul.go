package consul

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
)

type Registry struct {
	client *api.Client
}

// 透過單例模式建立 consul
var (
	consulClient *Registry
	once         sync.Once
	initErr      error
)

func New(consulAddr string) (*Registry, error) {
	once.Do(func() {
		config := api.DefaultConfig()
		config.Address = consulAddr
		client, err := api.NewClient(config)
		if err != nil {
			initErr = err
			return
		}
		consulClient = &Registry{
			client: client,
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return consulClient, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format")
	}
	port, _ := strconv.Atoi(parts[1])
	host := parts[0]
	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  false,
			TTL:                            "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceID, serviceName, hostPort string) error {
	logrus.WithFields(logrus.Fields{
		"instanceID":  instanceID,
		"serviceName": serviceName,
	}).Info("deregister from consul")
	return r.client.Agent().CheckDeregister(instanceID)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, e := range entries {
		ips = append(ips, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return ips, nil
}

func (r *Registry) HealthCheck(instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
