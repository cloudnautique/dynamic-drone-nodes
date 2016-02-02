package main

import (
	"github.com/drone/drone-go/drone"
	"github.com/rancher/go-rancher-metadata/metadata"
)

type RancherClient struct {
	client *metadata.Client
}

type DynamicNodePool interface {
	ListNodes(string) ([]*drone.Node, error)
}

func NewPoolClient(endpoint string) (DynamicNodePool, error) {
	client, err := NewRancherClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &RancherClient{client}, nil
}

func NewRancherClient(endpoint string) (*metadata.Client, error) {
	return metadata.NewClientAndWait(endpoint)
}

func (c *RancherClient) ListNodes(nodesPath string) ([]*drone.Node, error) {
	nodes := []*drone.Node{}

	rancherNodes, _ := c.client.GetServiceContainers("docker", "test")
	for _, rNode := range rancherNodes {
		n := &drone.Node{
			Addr: "tcp://" + rNode.PrimaryIp + ":2375",
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}
