package main

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/drone"
	"github.com/rancher/go-rancher-metadata/metadata"
)

type RancherClient struct {
	client *metadata.Client
}

type RancherStackServicePair struct {
	Stack   string
	Service string
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
	stackServicePair, err := segmentPathToStackServicePair(nodesPath)
	if err != nil {
		logrus.Fatal(err)
	}

	rancherNodes, err := c.client.GetServiceContainers(stackServicePair.Service, stackServicePair.Stack)
	if err != nil {
		logrus.Fatalf("Could not get containers from. %s", err)
	}

	for _, rNode := range rancherNodes {
		n := &drone.Node{
			Addr: "tcp://" + rNode.PrimaryIp + ":2375",
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// Path should be /stacks/<stackname>/services/<servicename>
func segmentPathToStackServicePair(path string) (*RancherStackServicePair, error) {
	rSSP := &RancherStackServicePair{}
	err := errors.New("Unable to parse path: " + path + "\n" +
		"expected format is: /stacks/<stackname>/services/<servicename>")

	path = strings.TrimPrefix(path, "/")
	segments := strings.Split(path, "/")
	if len(segments) != 4 {
		logrus.Error(len(segments))
		logrus.Error(segments)
		return rSSP, err
	}

	if segments[0] == "stacks" && segments[2] == "services" {
		rSSP.Stack = segments[1]
		rSSP.Service = segments[3]
		err = nil
	}

	return rSSP, err
}
