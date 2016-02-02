package main

import (
	"github.com/drone/drone-go/drone"
	"github.com/rancher/go-rancher-metadata/metadata"
)

type Client struct {
	EndpointURL string
}

type DynamicNodePool interface {
	ListNodes(string) ([]*drone.Node, error)
}

func NewPoolClient(endpoint string) *Client {
	return &Client{endpoint}
}

func (c *Client) ListNodes(nodesPath string) ([]*drone.Node, error) {
	nodes := []*drone.Node{}
	client, err := metadata.NewClientAndWait(c.EndpointURL)
	if err != nil {
		return nodes, err
	}

	rancherNodes, _ := client.GetServiceContainers("docker", "test")
	for _, rNode := range rancherNodes {
		n := &drone.Node{
			Addr: "tcp://" + rNode.PrimaryIp + ":2375",
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}
