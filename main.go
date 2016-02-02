package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/drone"
)

func main() {
	droneClient := &Drone{drone.NewClientToken("http://192.168.99.100:8000", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXh0IjoiY2xvdWRuYXV0aXF1ZSIsInR5cGUiOiJ1c2VyIn0.r_0gMglz2EwPkWISVQvCTF_0dIHG5QaNkZheyZqL7hc")}
	metadataClient := NewPoolClient("http://rancher-metadata/2015-12-19")

	serviceNodes, err := metadataClient.ListNodes("docker")
	if err != nil {
		logrus.Fatal(err)
	}

	if err = droneClient.ReconcileNodeDifferences(convertNodeToMap(serviceNodes)); err != nil {
		logrus.Fatal(err)
	}

}
