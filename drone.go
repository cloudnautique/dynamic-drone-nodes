package main

import (
	"fmt"
	"github.com/drone/drone-go/drone"
)

type NodeData struct {
	Node   *drone.Node
	Action string
}

type Drone struct {
	client drone.Client
}

func (d *Drone) getDroneNodeList() (map[string]*drone.Node, error) {
	nodes, err := d.client.NodeList()
	if err != nil {
		return nil, err
	}

	return convertNodeToMap(nodes), nil
}

// Gets the node list from Drone and compares it to the dynamic pool.
// It will take action to add / remove hosts
func (d *Drone) ReconcileNodeDifferences(poolView map[string]*drone.Node) error {
	actionableNodes := []*NodeData{}

	droneView, err := d.getDroneNodeList()
	if err != nil {
		return err
	}

	actionableNodes = append(actionableNodes, compareNodeMapsAndMark("ADD", poolView, droneView)...)
	actionableNodes = append(actionableNodes, compareNodeMapsAndMark("REMOVE", droneView, poolView)...)

	for _, node := range actionableNodes {
		fmt.Printf("%s: %s\n", node.Action, node.Node.Addr)
		switch node.Action {
		case "ADD":
			d.client.NodePost(node.Node)
		case "REMOVE":
			d.client.NodeDel(node.Node.ID)
		default:
			fmt.Errorf("Unknown Action")
		}
	}
	return nil
}

func compareNodeMapsAndMark(action string, listA, listB map[string]*drone.Node) []*NodeData {
	actionableNodes := []*NodeData{}

	// Loop through theh map and if A isn't in B
	// mark it for an action.
	for addrKey, node := range listA {
		if _, ok := listB[addrKey]; !ok {
			actionableNodes = append(actionableNodes, &NodeData{
				Node:   node,
				Action: action,
			})
		}
	}

	return actionableNodes
}

func convertNodeToMap(nodes []*drone.Node) map[string]*drone.Node {
	nodeMap := map[string]*drone.Node{}

	for _, node := range nodes {
		nodeMap[node.Addr] = node
	}

	return nodeMap
}
