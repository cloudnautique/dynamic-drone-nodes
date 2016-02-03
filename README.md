## Dynamic Drone Nodes

---- 

#### Purpose

This is a first pass at dynamically adding/removing nodes to a Drone CI server. The code watches Rancher metadata and adds/removes nodes from Drone CI server. The Rancher Metadata is considered authoritive, and this will reconcile to equal Rancher's view of the world.

##### Note

Right now it only handles 1 worker per address. Not necessarily a bad thing, just to add more workers you  have to add more containers making scheduling tricky.

#### Usage

```
NAME:
   dynamic-drone-nodes - Dynamically add and remove Drone CI nodes

USAGE:
   drone-dynamic-nodes [global options] command [command options] [arguments...]

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --drone-token 		API token for Drone CI [$DRONE_TOKEN]
   --drone-url 			URL for the Drone CI server [$DRONE_URL]
   --poll-interval "300"	Interval in (s) to poll dynamic pool
   --help, -h			show help
   --version, -v		print the version
```


