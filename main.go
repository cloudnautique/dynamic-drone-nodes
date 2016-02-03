package main

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cloudnautique/dynamic-drone-nodes/version"
	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

const (
	metadataURL = "http://rancher-metadata/2015-12-19"
)

func main() {
	app := cli.NewApp()
	app.Name = "dynamic-drone-nodes"
	app.Usage = "Dynamically add and remove Drone CI nodes"
	app.Action = appInit
	app.Version = version.VERSION
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
		cli.IntFlag{
			Name:  "poll-interval",
			Usage: "Interval in (s) to poll dynamic pool",
			Value: 300,
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

	if len(c.Args()) < 0 {
		logrus.Fatal("No path to watch found")
	}
	path := c.Args()[0]

	droneClient := &Drone{drone.NewClientToken(droneURL, droneToken)}

	poolClient, err := NewPoolClient(metadataURL)
	if err != nil {
		logrus.Fatal(err)
	}

	run(droneClient, poolClient, path, c.Int("poll-interval"))
}

func run(droneClient *Drone, poolClient DynamicNodePool, path string, interval int) {
	logrus.Infof("Path: %s", path)
	for {
		poolNodes, err := poolClient.ListNodes(path)
		if err != nil {
			logrus.Fatal(err)
		}

		if err = droneClient.ReconcileNodeDifferences(convertNodeToMap(poolNodes)); err != nil {
			logrus.Fatal(err)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
