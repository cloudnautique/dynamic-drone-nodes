package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

func main() {
	app := cli.NewApp()
	app.Name = "dynamic-drone-nodes"
	app.Usage = "Dynamically add and remove Drone CI nodes"
	app.Action = appInit
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "drone-token",
			Usage:  "API token for Drone CI",
			EnvVar: "DRONE_TOKEN",
		},
		cli.StringFlag{
			Name:   "drone-url",
			Usage:  "URL for the Drone CI server",
			EnvVar: "DRONE_URL",
		},
	}

	app.Run(os.Args)
}

func appInit(c *cli.Context) {
	droneToken := c.String("drone-token")
	if droneToken == "" {
		logrus.Fatal("Drone API token is missing")
	}

	droneURL := c.String("drone-url")
	if droneURL == "" {
		logrus.Fatal("Drone CI server URL is missing")
	}

	droneClient := &Drone{drone.NewClientToken(droneURL, droneToken)}

	PoolClient, err := NewPoolClient("http://rancher-metadata/2015-12-19")
	if err != nil {
		logrus.Fatal(err)
	}

	poolNodes, err := PoolClient.ListNodes("docker")
	if err != nil {
		logrus.Fatal(err)
	}

	if err = droneClient.ReconcileNodeDifferences(convertNodeToMap(poolNodes)); err != nil {
		logrus.Fatal(err)
	}
}
